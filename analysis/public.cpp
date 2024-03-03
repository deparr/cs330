#include <cstdio>
class C {
	int x = 42;
};
template <class M, class Secret>
struct public_cast {
	static inline M m{};
};

template <class Secret, auto M>
struct access {
	static const inline auto m
		= public_cast<decltype(M), Secret>::m = M;
};
template struct access<class CxSecret, &C::x>;

int main() {
	C c = C();
	int x = c.*public_cast<int C::*, CxSecret>::m;
	printf("c.x = %d\n", x);
}
