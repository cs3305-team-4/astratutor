import { Signalling, MESSAGE_TYPE } from "./signalling";
import { Peer } from "./peer";

// Adapted From:
// https://w3c.github.io/webrtc-pc/#perfect-negotiation-example
export class WebRTCHandler {
    signalling: Signalling
    media: { [id: string]: MediaStream }
    peers: { [id: string]: Peer }

    // Callbacks
    ontrack?: (id: string, event: RTCTrackEvent) => any

    constructor(signalling: Signalling, media?: MediaStream[]) {
        this.signalling = signalling
        this.peers = {}
        this.media = {}

        if (!media) return
        media.forEach(stream => {
            stream.getTracks().forEach(track => {
                this.media[track.id] = stream
            })
        })
    }

    addPeer(id: string, polite?: boolean): Peer {
        console.log("Adding Peer: " + id)
        //TODO(james): STUN and TURN server details
        const peer = new Peer(new RTCPeerConnection(), polite ?? false)

        Object.entries(this.media).forEach(([id, stream]) => {
            const sender = peer.conn.addTrack(stream.getTrackById(id)!, stream)
            peer.tracks[id] = sender
        })

        if (Object.keys(this.media).length == 0 && !polite) {
            console.log("No media. Sending dummy to start session")
            peer.conn.addTransceiver("audio")
        }

        peer.conn.onicecandidate = (event) => {
            if (event.candidate) {
                console.log("Sending Candidate: " + id)
                this.signalling.send(MESSAGE_TYPE.CANDIDATE, id, event.candidate)
            }
        }

        peer.conn.onnegotiationneeded = async () => {
            console.log("Negotiaion Needed: " + id)
            try {
                peer.makingOffer = true
                await peer.conn.setLocalDescription({})
                this.signalling.send(MESSAGE_TYPE.SDP, id, peer.conn.localDescription)
            } catch (err) {
                console.error(err)
            } finally {
                peer.makingOffer = false
            }
        }

        peer.conn.ontrack = (event) => {
            if (!this.ontrack) return
            this.ontrack(id, event)
        }

        this.peers[id] = peer
        return peer
    }

    async incomingSDP(id: string, sdp: RTCSessionDescription) {
        const peer = this.peers[id] ?? this.addPeer(id, true)
        try {
            const readyForOffer =
                !peer.makingOffer &&
                (peer.conn.signalingState == "stable" || peer.isSettingRemoteAnswerPending)
            const offerCollision = sdp.type === "offer" && !readyForOffer
            peer.ignoreOffer = !peer.polite && offerCollision
            if (peer.ignoreOffer) return

            peer.isSettingRemoteAnswerPending = sdp.type === "answer"
            await peer.conn.setRemoteDescription(sdp)
            peer.isSettingRemoteAnswerPending = false

            if (sdp.type === "offer") {
                await peer.conn.setLocalDescription({})
                this.signalling.send(MESSAGE_TYPE.SDP, id, peer.conn.localDescription)
            }
        } catch (err) {
            console.error(err)
        }
    }

    async incomingCandidate(id: string, candidate: RTCIceCandidate) {
        const peer = this.peers[id]
        try {
            try {
                await peer.conn.addIceCandidate(candidate)
            } catch (err) {
                if (!peer.ignoreOffer) throw err;
            }
        } catch (err) {
            console.error(err)
        }
    }

    addTrack(track: MediaStreamTrack, stream?: MediaStream) {
        console.log("Adding Track: " + track.id)
        stream = stream ?? new MediaStream([track])
        this.media[track.id] = stream
        Object.values(this.peers).forEach(peer => {
            const sender = peer.conn.addTrack(track, stream ?? new MediaStream([track]))
            peer.tracks[track.id] = sender
        })
    }

    removeTrack(track: MediaStreamTrack) {
        if (!this.media[track.id]) return
        this.media[track.id].removeTrack(track)
        delete this.media[track.id]
        Object.values(this.peers).forEach(peer => {
            peer.conn.removeTrack(peer.tracks[track.id])
            delete peer.tracks[track.id]
        })
    }

    removeTrackById(id: string) {
        if (!this.media[id]) return
        this.removeTrack(this.media[id].getTrackById(id)!)
    }

    replaceTrackById(oldid: string, track: MediaStreamTrack) {
        Object.values(this.peers).forEach(peer => {
            peer.tracks[oldid].replaceTrack(track)
            peer.tracks[track.id] = peer.tracks[oldid]
            delete peer.tracks[oldid]
        })
    }
}