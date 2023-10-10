import { io } from "socket.io-client";

const nestBaseUrl: string = process.env.NEXT_PUBLIC_NEST_BASE_URL ?? ''

export const socket = io(nestBaseUrl, {
    autoConnect: false,
})