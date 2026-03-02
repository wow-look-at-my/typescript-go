//// [tests/cases/compiler/declarationEmitCommonSourceDirectoryDoesNotContainAllFiles.ts] ////

//// [index.ts]
export * from "./src/"
//// [index.ts]
export class B {}
//// [index.ts]
import { B } from "b";

export default function () {
	return new B();
}
//// [index.ts]
export * from "./src/"



//// [index.d.ts]
import { B } from "b";
export default function () {
    return new B();
}
export * from "./src/";
