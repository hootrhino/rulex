#include <stdio.h>
#include <stdlib.h>
#include <time.h>
int main(int argc, char const *argv[])
{
    FILE *fp = NULL;
    fp = fopen("serial_number.csv", "w");
    if (fp == NULL)
    {
        printf("file can't be opened\n");
        exit(1);
        return 0;
    }
    srand((unsigned)time(NULL));
    fprintf(fp, "Address , Channel\n");
    for (size_t i = 1; i < 0xFF; i++)
    {
        fprintf(fp, "%d , %d\n", i, rand() % 0xFFFF + 10000);
    }
    fclose(fp);
    return 0;
}
