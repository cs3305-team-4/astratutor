export interface Callbacks {
  onopen?: (ev: Event) => any;
  onmessage?: (ev: MessageEvent) => any;
  onclose?: (ev: CloseEvent) => any;
  onerror?: (ev: Event) => any;
}

export enum MESSAGE_TYPE {
  AHOY_HOY = 1,
  CHAT,
  SDP,
  CANDIDATE,
}

export class Signalling {
  ws: WebSocket;
  id: string;

  //TODO(james): Better callback system
  constructor(id: string, classUrl: string, callbacks: Callbacks) {
    this.id = id;

    this.ws = new WebSocket(classUrl);
    this.ws.onopen = callbacks.onopen!;
    this.ws.onmessage = callbacks.onmessage!;
    this.ws.onclose = callbacks.onerror!;
    this.ws.onerror = callbacks.onerror!;
  }

  send(message_type: MESSAGE_TYPE, to: string, data: any): void {
    console.log(this.ws.readyState);
    this.ws.send(JSON.stringify({ src: this.id, dest: to, type: message_type, data: data }));
  }
}
