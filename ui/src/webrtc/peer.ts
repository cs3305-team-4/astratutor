import { StreamType } from './stream_types';

export class Peer {
  id: string;
  ready: boolean;
  conn: RTCPeerConnection;
  polite: boolean;
  makingOffer: boolean;
  ignoreOffer: boolean;
  isSettingRemoteAnswerPending: boolean;
  tracks: { [id: string]: RTCRtpSender };
  correlateChan?: RTCDataChannel;
  streamCorrelations: { [sid: string]: StreamType };

  constructor(id: string, conn: RTCPeerConnection, polite: boolean) {
    this.id = id;
    this.ready = false;
    this.conn = conn;
    this.polite = polite || false;
    this.makingOffer = false;
    this.ignoreOffer = false;
    this.correlateChan = undefined;
    this.tracks = {};
    this.isSettingRemoteAnswerPending = false;
    this.streamCorrelations = {};
  }
}
