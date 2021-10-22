#include <stdio.h>
#include <string.h>
typedef struct __attribute__((__packed__))
{
    unsigned int sw0 : 1;
    unsigned int sw1 : 1;
    unsigned int sw2 : 1;
    unsigned int sw3 : 1;
    unsigned int sw4 : 1;
    unsigned int sw5 : 1;
    unsigned int sw6 : 1;
    unsigned int sw7 : 1;
    unsigned int sw8 : 1;
    unsigned int sw9 : 1;
} modbus_data;

void write_data()
{
    FILE *fp = fopen("modbus_data.bin", "wb");
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
    printf("Write: %d %d %d %d %d %d %d %d %d %d\n",
           data.sw0,
           data.sw1,
           data.sw2,
           data.sw3,
           data.sw4,
           data.sw5,
           data.sw6,
           data.sw7,
           data.sw8,
           data.sw9);
    fwrite(&data, sizeof(modbus_data), 1, fp);
    fclose(fp);
}
void read_data()

{
    FILE *fp = fopen("modbus_data.bin", "rb");
    modbus_data data;
    fread(&data, sizeof(modbus_data), 1, fp);
    printf("Read: %d %d %d %d %d %d %d %d %d %d\n",
           data.sw0,
           data.sw1,
           data.sw2,
           data.sw3,
           data.sw4,
           data.sw5,
           data.sw6,
           data.sw7,
           data.sw8,
           data.sw9);
    fclose(fp);
}

int main(int argc, char const *argv[])
{
    write_data();
    read_data();
}