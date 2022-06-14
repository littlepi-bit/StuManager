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

```c++
		uint8_t *result = DHT11(vol,mode);
		float distance = HCSR04_GetDistant(&htim3,GPIOG,GPIO_PIN_2,GPIOG,GPIO_PIN_3);
		printf("Distance:%lf\n",distance);
		int t = 103*vol/4095;
		t %= 101;
		if(mode == 0){
			printf("vol:%d Voltage:%.3fv\r\n",vol,(double)vol*3.3/(double)4095);
			__HAL_TIM_SetCompare(&htim1,TIM_CHANNEL_1,516-vol*500/4095);
			for(int i = 0;i < 300;i++){
				int tmp = t;
				if (tmp == 0){
					ChoseW(0);
					ShowNum(0);
					for(int k = 0;k < 5000;k++);
				}
				int j = 0;
				while(tmp>0){
					ChoseW(j++);
					ShowNum(tmp%10);
					tmp /= 10;
					for(int k = 0;k < 5000;k++);
				}
			}
		}
		else if (mode == 1){
			for(int i = 0;i < 120;i++){
				uint8_t t0 = result[3];
				ChoseW(0);
				ShowNum(12);
				for(int k = 0;k < 8000;k++);
				int j = 1;
				ChoseW(j++);
				ShowNum(t0%10);
				for(int k = 0;k < 8000;k++);
				uint8_t t1 = result[2];
				ChoseW(j++);
				ShowNumPoint(t1%10);
				for(int k = 0;k < 8000;k++);
				t1/=10;
				ChoseW(j++);
				ShowNum(t1%10);
				for(int k = 0;k < 8000;k++);
				t1/=10;
				ChoseW(j++);
				ShowNum(19);
				for(int k = 0;k < 8000;k++);
				uint8_t t2 = result[1];
				ChoseW(j++);
				ShowNum(t2%10);
				for(int k = 0;k < 8000;k++);
				uint8_t t3 = result[0];
				ChoseW(j++);
				ShowNumPoint(t3%10);
				for(int k = 0;k < 8000;k++);
				t3/=10;
				ChoseW(j++);
				ShowNum(t3%10);
				for(int k = 0;k < 8000;k++);
			}
		}
```