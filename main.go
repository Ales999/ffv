package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/alecthomas/kong"
	"github.com/ales999/cisaccs"
	"github.com/ales999/ffv/utils"
)

var cli struct {
	// Обязательный аргумент - список cisco хостов к которым будем подключаться
	CheckHosts []string `arg:"" name:"hosts" help:"Name of cisco hosts for finded IP"`
	// Флаг на уникальность вывода
	UniqueOutput bool `help:"Вывод будет один общий" short:"u" default:"false"`
	// Номер порта для ssh
	PortSsh int `help:"SSH порт для доступа к cisco" short:"p" default:"22"`
	// Путь к файлу конфигурации имя_cisco/группа/ip - env: CISFILE
	CisFileName string `help:"Путь к файлу конфигурации имя_cisco/группа/ip" default:"/etc/cisco/cis.yaml" env:"CISFILE"`
	// Путь к файлу конфигурации имя_группы/имя/пароль - env: CISPWDS
	PwdFileName string `help:"Путь к файлу конфигурации имя_группы/имя/пароль" default:"/etc/cisco/passw.json" env:"CISPWDS"`
}

//var skipVlans []string

func main() {

	ctx := kong.Parse(&cli,
		kong.Name("ffv"),
		kong.Description("Find Free Vlans"),
		kong.UsageOnError(),
	)

	if currentUserUid == "0" {
		fmt.Println("Запрет запуска под root")
		os.Exit(1)
	}

	err := findFreeVlan(cli.CheckHosts)
	ctx.FatalIfErrorf(err)
	os.Exit(0)

}

func findFreeVlan(hosts []string) error {

	var vlans []utils.VlanLineData
	// Что будем выполнять на cisco
	cmds := []string{"sh vlan br"}

	// Подготовка к подключению.
	acc := cisaccs.NewCisAccount(cli.CisFileName, cli.PwdFileName)

	// Количество хостов в списке на получение данных
	hstcount := len(hosts)

	// Пробегаем по всем хостам
	for _, hst := range hosts {
		//vlans = []string{} // Масив vlan-ов
		// Получаем данные с каждого хоста подключаясь к каждому по очереди.
		cisout, err := acc.OneCisExecuteSsh(hst, cli.PortSsh, cmds)
		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		// Парсим чего выдали нам cisco
		for _, line := range cisout {
			// Парсим строку
			_nvl := utils.ParseVlan(line)
			// Проверка что запись не пустая
			if _nvl.GetId() == 0 {
				continue
			}
			// Если строка не пуста - добавляем
			vlans = append(vlans, _nvl)
		}
		// Если нет уникального вывода то печатаем все подрят, по одному хосту.
		if !cli.UniqueOutput && vlans != nil {
			// Выполняем поиск свободных vlal-ов
			fr := GenerateRange(&vlans)
			// Если хост один.
			if hstcount == 1 {
				// Печать результата
				PrintFreeRange(&fr)
			} else {
				// Выводим имя хоста
				fmt.Println("----------")
				fmt.Println(hst, ":")
				// Печать результата
				PrintFreeRange(&fr)
			}
			// Удалаяем все из среза
			vlans = nil
			vlans = []utils.VlanLineData{}
		}

	} // End for hosts
	// Если список не пустой, значит нам нужнен уникальный список по всем хостам
	if cli.UniqueOutput && vlans != nil {
		// Выполняем поиск свободных vlal-ов в общем списке
		fr := GenerateRange(&vlans)
		// Печать результата
		PrintFreeRange(&fr)
	}

	return nil
}

// Печать результатов
func PrintFreeRange(freerange *[]FreeRange) {

	for _, v := range *freerange {
		v.PrintData()
	}

}

func GenerateRange(vlans *[]utils.VlanLineData) []FreeRange {

	// срез vlan-ов что уже заняты на коммутаторе
	var zanvls []int
	// Последний разрешенный номер VLAN-а
	const LASTVLANID int = 4095
	var fr FreeRange      // Тут будем хранить начало и стар своодного диапазона
	var fouts []FreeRange // Список всех диапазонов

	// Заполняем его.
	for _, vl := range *vlans {
		zanvls = append(zanvls, vl.GetId())
	}
	// Если нам нужна уникальность то в списке занятых - zanvls, все подрят и надо созать лбщий список.
	if cli.UniqueOutput {
		// Сортируем
		sort.Ints(zanvls)
		// и Удаляем дибликаты
		zanvls = RemoveDuplicateInt(zanvls)
	}

	// Сохраним размер дабы не выйти за границу диапазона далее
	var zlen = len(zanvls)

	// // Основной ЦИКЛ по номерам VLAN-ов которые могут быть в системе (1-й всегда есть).

	for zid := range zanvls {
		var _start int
		var _stop int

		var _rand int
		if (zid + 1) < zlen {
			_rand = zanvls[zid+1] - zanvls[zid]
			if _rand > 1 {
				_start = zanvls[zid] + 1
				_stop = zanvls[zid+1] - 1
			}
		} else {
			_rand = LASTVLANID - zanvls[zid]
			if _rand > 1 {
				_start = zanvls[zid] + 1
				_stop = LASTVLANID
			}
		} // endif
		// Если на этом этапе цикла найден свободный диапазон - добавим его
		if _rand > 1 {
			fr = *NewFreeRange(_start, _stop)
			fouts = append(fouts, fr)
		}

	} // end for

	/*
		// Test print data
		for _, v := range fouts {
			v.PrintData()
		}
	*/

	return fouts

}
