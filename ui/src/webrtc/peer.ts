export class Peer {
    conn: RTCPeerConnection
    polite: boolean
    makingOffer: boolean
    ignoreOffer: boolean
    isSettingRemoteAnswerPending: boolean
    tracks: { [id: string]: RTCRtpSender}

    constructor(conn: RTCPeerConnection, polite: boolean) {
        this.conn = conn
        this.polite = polite ?? false
        this.makingOffer = false
        this.ignoreOffer = false
        this.isSettingRemoteAnswerPending = false
        this.tracks = {}
    }
}