#include <stdio.h>
#include <string.h>
typedef struct
{
    int sw0 : 1;
    int sw1 : 1;
    int sw2 : 1;
    int sw3 : 1;
    int sw4 : 1;
    int sw5 : 1;
    int sw6 : 1;
    int sw7 : 1;
    int sw8 : 1;
    int sw9 : 1;
} modbus_data;

int main(int argc, char const *argv[])
{
    FILE *p = fopen("modbus_data.bin.", "wb");
    modbus_data data;
    data.sw0 = 1;
    data.sw1 = 0;
    data.sw2 = 1;
    data.sw3 = 1;
    data.sw4 = 0;
    data.sw5 = 1;
    data.sw6 = 1;
    data.sw7 = 0;
    data.sw8 = 0;
    data.sw9 = 1;
    fwrite(&data, sizeof(modbus_data), 1, p);
    fclose(p);
    return 0;
}
