```c++
printf("vol:%d Voltage:%.3fv\r\n",vol,(double)vol*3.3/(double)4095);
		__HAL_TIM_SetCompare(&htim1,TIM_CHANNEL_1,516-vol*500/4095);
		int t = 103*vol/4095;
		for(int i = 0;i < 20;i++){
			int tmp = t;
			if (tmp == 0){
				ChoseW(0);
				ShowNum(0);
				for(int k = 0;k < 10000;k++);
			}
			int j = 0;
			while(tmp>0){
				ChoseW(j++);
				ShowNum(tmp%10);
				tmp /= 10;
				for(int k = 0;k < 10000;k++);
			}
		}
```

```c++
uint8_t ch;
uint8_t ch_r;
int fputc(int c, FILE * f)
{
	ch=c;
	HAL_UART_Transmit(&huart1,&ch,1,1000);
	return c;
}
int fgetc(FILE * F) 
{
	HAL_UART_Receive (&huart1,&ch_r,1,0xffff);
	return ch_r;
}

```