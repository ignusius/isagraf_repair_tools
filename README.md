Программа испраляет пробоемы с G2 и G3 ошибками в IsaGRAF 3.x, когда виpуально всё должно работать.

Error: MAX @VAR,@BOX,@TXT: ID > 64785 - Превышен лемнит на ID для блоков.

Исправить эту ошибку можно либо удалением записи в .lsf-файле или изменением на id который не используется.

Пример некорректного id
```
@BOX:64656,Y=N,P=(19,283),S=(1,1),C=(1,1),X={\div}
@BOX:64658,Y=N,P=(2,286),S=(1,1),C=(1,1),X={\div}
@TXT:64666,Y=C32768,P=(35,320),S=(10,2),C=(0,0),X=25{
         Газ в автомате  
}
@VAR:65417,Y=N,P=(69,421),S=(8,1),C=(1,1),X=Gaz_Tm8  <------- Плохой ID > 64785
@ARC:1,D=2,Y=N,Z=(50,50),F=(127,0),T=(126,1)
```
