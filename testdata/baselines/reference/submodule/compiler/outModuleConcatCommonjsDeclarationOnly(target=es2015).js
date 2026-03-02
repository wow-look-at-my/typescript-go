//// [tests/cases/compiler/outModuleConcatCommonjsDeclarationOnly.ts] ////

//// [a.ts]
export class A { }

//// [b.ts]
import {A} from "./ref/a";
export class B extends A { }


//// [all.js]
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.A = void 0;
class A {
}
exports.A = A;
Object.defineProperty(exports, "__esModule", { value: true });
exports.B = void 0;
const a_1 = require("./ref/a");
class B extends a_1.A {
}
exports.B = B;
//# sourceMappingURL=all.js.map