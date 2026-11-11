import { TextDecoder, TextEncoder } from "util";

// react-router v7 uses TextEncoder, which the Jest 27 jsdom environment lacks
Object.assign(globalThis, { TextDecoder, TextEncoder });
