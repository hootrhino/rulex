#include <stdio.h>
typedef int (*func)(int);
int run(func funcs[], int len, int arg);
int pipline(int acc, func funcs[], int len, int arg);

int add1(int a)
{
    return a + 1;
}
int sub1(int a)
{
    return a - 1;
}

int main(int argc, char const *argv[])
{
    func funcs[5];
    funcs[0] = add1;
    funcs[1] = sub1;
    funcs[2] = sub1;
    funcs[3] = add1;
    funcs[4] = add1;
    printf("pipline result: %d\n", run(funcs, 5, 100));
    return 0;
}

int run(func funcs[], int len, int arg)
{
    return pipline(0, funcs, len, arg);
}
int pipline(int acc, func funcs[], int len, int arg)
{
    if ((acc) == (len - 1))
    {
        return funcs[acc](arg);
    }
    else
    {
        return pipline(acc + 1, funcs, len, funcs[acc](arg));
    }
}
