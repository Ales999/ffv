# ffv

**F**ind **F**ree **V**lan for cisco commutator

Example run:

```bash
./ffv cisco3650
2-9
11-28
30
32 - 1001
1006 - 1999
2001 - 4095
```

This prog used my [cisaccs](https://github.com/Ales999/cisaccs) library for access to cisco by ssh.

Программа для поиска свободных диапазонов VLAN-ов на коммутаторах cisco.
Можно задать через пробел сразу несколько имен коммутаторов.

Добавлены флаги:
 * -r - подсчет используемых
 * -u - если коммутаторов несколько то только уникальные выводим, в виде summary по всем.
 * -с - вывести общее количество в конце списка используемых или свободных VLAN-ов
