#include "login.h"
#include <QApplication>

int main(int argc, char *argv[]){
    QApplication a(argc, argv);

    Login w;
    w.show();
    return a.exec();
}


//#include <QDebug>
//#include <functional>

//class A
//{
//public :
//    A()
//    {
//        qDebug() << "A " << this;
//    }
//    A(const A&a)
//    {
//        qDebug() << "A(&a)" << this << &a;
//        this->value = a.value;
//    }

//    ~A()
//    {
//        qDebug() << "~A"<< this;
//    }
//public :
//    int value;
//};

//std::function<void ()> retFunction()
//{
//    A a;
//    a.value = 14;
//    qDebug() << "1";
//    return [=]() {
//        qDebug() << "2";
//        qDebug() << a.value;
//    };
//}

//int main() {
//    {
//        auto f1 = retFunction();
//        qDebug() << "3";
//        f1();
//        qDebug() << "5";
//        f1();
//        qDebug() << "6";
//    }
//    qDebug() << "4";
//    qDebug() << "xxxx";
//}
