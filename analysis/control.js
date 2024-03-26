let tag = 0;
function T(l, v, r) {
	tag++;
    return {left: l, value: v, right: r, tag}
}


let e = null;
let t1 = T(T(T(e, 1, e), 2, T(e, 3, e)), 4, T(e, 5, e));
let t2 = T(e, 1, T(e, 2, T(e, 3, T(e, 4, T(e, 5, e)))));
let t3 = T(e, 1, T(e, 2, T(e, 3, T(e, 4, T(e, 6, e)))));

console.log(t1);
// console.log(t2);
console.log(t3);



function coroutine_wrap(f) {
	async function yield(val) {

		return 
	}

	return () => f(yield);
}

/**
 * @param {Promise<co>} co coroutine promise
 * @returns {Promise<co>} 
 *
*/
function resume(co) {
	return co.then((val) => val)

}

async function visit(t, yield) {
	if (t) {
		visit(t.left)
		await yield(t.value)
		console.log(t.tag, t.value)
		visit(t.right)
	}
}

function cmp() {
	let co1 = coroutine_wrap((yield) => visit(t1, yield))
	let co2 = coroutine_wrap((yield) => visit(t3, yield))

	while (true) {
		let v1 = resume(co1)
		let v2 = resume(co2)

		if (v1 != v2) {
			return false
		}
		if (v1 == null && v2 == null) {
			return true
		}
	}
}

console.log(cmp(t1, t3))
