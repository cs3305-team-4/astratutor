export interface Callbacks {
  onopen?: (ev: Event) => void;
  onmessage?: (ev: MessageEvent) => void;
  onclose?: (ev: CloseEvent) => void;
  onerror?: (ev: Event) => void;
}

export enum MESSAGE_TYPE {
  AHOY_HOY = 1,
  CHAT,
  SDP,
  CANDIDATE,
  STOP_STREAM,
  DRAW,
  WIPE,
  UNDO,
  CHANGE_BG,
  INIT,
  LEAVE,
}

export class Signalling {
  ws: WebSocket;
  id: string;

  constructor(id: string, classUrl: string, callbacks: Callbacks) {
    this.id = id;

    this.ws = new WebSocket(classUrl);
    if (callbacks.onopen) this.ws.onopen = callbacks.onopen;
    if (callbacks.onmessage) this.ws.onmessage = callbacks.onmessage;
    if (callbacks.onclose) this.ws.onclose = callbacks.onclose;
    if (callbacks.onerror) this.ws.onerror = callbacks.onerror;
  }

  send(message_type: MESSAGE_TYPE, to: string, data: any): void {
    console.log('Sending: ' + this.ws.readyState + ' - ' + message_type);
    this.ws.send(JSON.stringify({ src: this.id, dest: to, type: message_type, data: data }));
  }

  onopen(func: (ev: Event) => void) {
    this.ws.onopen = func;
  }

  onmessage(func: (ev: MessageEvent) => void) {
    this.ws.onmessage = func;
  }

  onclose(func: (ev: CloseEvent) => void) {
    this.ws.onclose = func;
  }

  onerror(func: (ev: Event) => void) {
    this.ws.onerror = func;
  }
}
