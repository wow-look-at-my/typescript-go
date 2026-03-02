//// [tests/cases/conformance/types/uniqueSymbol/uniqueSymbolsDeclarationsInJsErrors.ts] ////

//// [uniqueSymbolsDeclarationsInJsErrors.js]
class C {
    /**
     * @type {unique symbol}
     */
    static readwriteStaticType;
    /**
     * @type {unique symbol}
     * @readonly
     */
    static readonlyType;
    /**
     * @type {unique symbol}
     */
    static readwriteType;
}


//// [uniqueSymbolsDeclarationsInJsErrors-out.js]
"use strict";
class C {
    /**
     * @type {unique symbol}
     */
    static readwriteStaticType;
    /**
     * @type {unique symbol}
     * @readonly
     */
    static readonlyType;
    /**
     * @type {unique symbol}
     */
    static readwriteType;
}
