//// [tests/cases/compiler/outModuleConcatES6.ts] ////

//// [a.ts]
export class A { }

//// [b.ts]
import {A} from "./ref/a";
export class B extends A { }

//// [all.js]
export class A {
}
import { A } from "./ref/a";
export class B extends A {
}
//# sourceMappingURL=all.js.map