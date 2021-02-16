import {
  AudioOutlined,
  CameraFilled,
  DesktopOutlined,
  PhoneFilled,
  SettingFilled,
  VideoCameraOutlined,
} from '@ant-design/icons';
import { Button, Col, Layout, Modal, Row, Select, Tooltip, Typography } from 'antd';
import React, { ReactElement, useContext, useEffect, useRef } from 'react';
import { useAsync } from 'react-async-hook';
import { useHistory, useParams } from 'react-router-dom';
import styled from 'styled-components';
import { APIContext } from '../api/api';
import { AccountType, ProfileResponseDTO } from '../api/definitions';
import Messaging, { Message } from '../components/Messaging';
import { UserAvatar } from '../components/UserAvatar';
import { SettingsCTX } from '../api/classroom';
import { Signalling, MESSAGE_TYPE } from '../webrtc/signalling';
import { WebRTCHandler } from '../webrtc/webrtc';
import { StreamType } from '../webrtc/stream_types';
import { screenStream } from '../webrtc/devices';

interface IWebcam {
  profile: ProfileResponseDTO;
  ref: React.ReactElement<HTMLVideoElement>;
  stream: MediaStream;
  streaming: boolean;
}

const webcamHeight = 200;

const StyledLayout = styled(Layout)`
  background-color: rgb(21 20 20);
  color: #fff;
  position: fixed;
  top: 0;
  left: 0;
  width: 100vw;
  height: 100vh;
  z-index: 100;
  display: flex;
  justify-content: center;
`;

const StyledSider = styled(Layout.Sider)`
  background-color: rgb(15 15 15);
  border-right: 2px solid rgb(10 10 10);
`;

const StyledWebcam = styled.div<{ index: number }>`
  width: 100%;
  padding: 0;
  margin: 0;
  overflow: hidden;
  height: ${webcamHeight}px;
  & .ant-typography {
    padding-left: 10px;
    color: #fff;
  }
  & video {
    width: 100%;
  }
  & .profile {
    background: linear-gradient(180deg, rgba(255, 255, 255, 0) 0%, rgba(0, 0, 0, 1) 100%);
    padding: 5px 10px;
    height: 42px;
    position: absolute;
    width: 100%;
    top: ${(props) => props.index * webcamHeight + 158}px;
  }
`;

const StyledIcon = styled(SettingFilled)`
  color: rgb(192, 192, 192);
  font-size: 30px;
  position: fixed;
  top: 70px;
  right: 15px;
  cursor: pointer;
`;

const StyledTools = styled(Layout.Footer)`
  background-color: rgb(5 5 5);
  position: absolute;
  bottom: 0;
  left: 300px;
  width: calc(100% - 300px);
  text-align: center;
  vertical-align: middle;
`;

const StyledVideo = styled.video`
  background-color: #000;
  width: 100%;
  height: calc(100% - 88px);
`;

export function LessonClassroom(): ReactElement {
  const { lid } = useParams<{ lid: string }>();
  const settings = useContext(SettingsCTX);
  const history = useHistory();
  const api = useContext(APIContext);

  const [messages, setMessages] = React.useState<Message[]>([]);
  const [webcamDisplays, setWebcamDisplays] = React.useState<IWebcam[]>([]);
  const [settingsOpen, setSettingsOpen] = React.useState(false);
  const [webcamEnabled, setWebcamEnabled] = React.useState(true);
  const [screenEnabled, setScreenEnabled] = React.useState(false);
  const [micEnabled, setMicEnabled] = React.useState(true);
  const [screen, setScreen] = React.useState<MediaStream>();

  const signalling = settings.signalling;
  const handler = useRef<WebRTCHandler>();
  const [addingPeer, setAddingPeer] = React.useState(false);
  const screenRef = useRef<HTMLVideoElement>();

  useAsync(async () => {
    // Signalling can be none if classroom page is refreshed before being sent back to lobby
    if (signalling == null) return;
    handler.current = new WebRTCHandler(signalling, () => setAddingPeer(true));
    console.log(handler.current);
    signalling.onmessage(async (event: MessageEvent) => {
      const message = JSON.parse(event.data);
      // Should never happen but just in case
      if (message.src === api.claims?.sub) return;

      const type: MESSAGE_TYPE = message.type;
      switch (type) {
        case MESSAGE_TYPE.AHOY_HOY: {
          handler.current?.addPeer(message.src);
          break;
        }
        case MESSAGE_TYPE.CHAT: {
          console.log('New Message: ', message);
          const messageData = Object.assign(message.data, { date: new Date(message.data.date) });
          setMessages((prev) => prev.concat(messageData));
          break;
        }
        case MESSAGE_TYPE.SDP: {
          // Only respond to SDP destined for us
          if (message.dest !== signalling.id) return;
          await handler.current?.incomingSDP(message.src, message.data);
          break;
        }
        case MESSAGE_TYPE.CANDIDATE: {
          // Only respond to candidates destined for us
          if (message.dest !== signalling.id) return;
          await handler.current?.incomingCandidate(message.src, message.data);
          break;
        }
      }
    });
    signalling.send(MESSAGE_TYPE.AHOY_HOY, '', null);

    handler.current.ontrack = (id: string, correlation: StreamType, event: RTCTrackEvent) => {
      console.log('NEW TRACK', id, correlation, event);
      switch (correlation) {
        case StreamType.Camera:
          setWebcamDisplays((prev) => {
            const other = prev.findIndex((v) => v.profile.account_id === id);
            console.log(webcamDisplays, other, id);
            if (other > -1) {
              prev.splice(other, 1);
            }
            return prev;
          });
          const stream = event.streams.length ? event.streams[0] : new MediaStream();
          const ref = (
            <video
              key={id}
              ref={(ref) => {
                if (ref) {
                  ref.srcObject = stream;
                  ref.play();
                }
              }}
            />
          );
          const profile = settings.otherProfiles[id];
          console.log(settings.otherProfiles);
          setWebcamDisplays((prev) => prev.concat({ stream, ref, profile, streaming: true }));
          break;
      }
    };
  }, []);

  useAsync(async () => {
    const last = messages[messages.length - 1];
    if (last && !last.profile) {
      console.log(last);
      const profile = await api.services.readProfileByAccountID(
        api.account?.id ?? '',
        api.account?.type ?? AccountType.Tutor,
      );
      signalling?.send(MESSAGE_TYPE.CHAT, '', { text: last.text, date: last.date, profile });
    }
  }, [api.account?.id, messages]);

  const addWebcam = (web: IWebcam) => {
    console.log('Adding Webcam to', handler.current?.peers);
    const other = webcamDisplays.findIndex((v) => v.ref.key === web.ref.key);
    if (other > -1) {
      console.log(other);
      setWebcamDisplays((prev) => {
        const temp = prev;
        const tracks = web.stream.getTracks();
        const old = temp[other];
        if (old.streaming) {
          tracks.forEach((v, i) => {
            handler.current?.replaceTrack(temp[other].stream.getTracks()[i], v);
          });
        } else {
          web.stream.getTracks().forEach((v) => {
            handler.current?.addTrack(v, StreamType.Camera, web.stream);
          });
        }
        delete temp[other];
        temp[other] = web;
        return temp;
      });
    } else {
      web.streaming = Object.keys(handler.current!.peers).length > 0;
      setWebcamDisplays((prev) => prev.concat(web));
      web.stream.getTracks().forEach((v) => {
        handler.current?.addTrack(v, StreamType.Camera, web.stream);
      });
    }
  };

  useAsync(async () => {
    // Get Self Profile
    if (!webcamEnabled) {
      settings.webcamStream?.getVideoTracks().forEach((v) => {
        v.enabled = false;
      });
    } else {
      settings.webcamStream?.getVideoTracks().forEach((v) => {
        v.enabled = true;
      });
    }
    const video = (
      <video
        key={'self'}
        ref={(ref) => {
          if (ref) {
            ref.srcObject = webcamEnabled ? settings.webcamStream : null;
            webcamEnabled && ref.play();
            navigator.mediaDevices.getSupportedConstraints();
          }
        }}
      />
    );
    if (api.account && settings.webcamStream) {
      const web: IWebcam = {
        profile: await api.services.readProfileByAccount(api.account),
        ref: video,
        stream: settings.webcamStream,
        streaming: false,
      };
      addWebcam(web);
      if (addingPeer) {
        setAddingPeer(false);
      }
    }
  }, [settings.webcamStream, webcamEnabled, addingPeer]);

  useAsync(async () => {
    if (screenEnabled) {
      const src = await screenStream();
      if (!src) {
        setScreenEnabled(false);
        return;
      }
      for (const track of src.getTracks()) {
        handler.current?.addTrack(track, StreamType.Screen, src);
      }
      setScreen(src);
      if (screenRef.current) {
        screenRef.current.srcObject = src;
      }
    } else {
      screen?.getTracks().forEach((v) => {
        v.enabled = false;
        handler.current?.removeTrack(v);
      });
      setScreen(undefined);
    }
  }, [screenEnabled]);

  const hangup = () => {
    settings.webcamStream?.getVideoTracks().forEach((v) => {
      v.stop();
    });
    history.push(`/lessons/${lid}/goodbye`);
  };

  return (
    <StyledLayout>
      <StyledLayout>
        <StyledIcon style={{ zIndex: 1000 }} onClick={() => setSettingsOpen(true)} title="Settings" />
        <Modal
          title="Settings"
          visible={settingsOpen}
          centered
          onOk={() => setSettingsOpen(false)}
          onCancel={() => setSettingsOpen(false)}
        >
          <Layout>
            <br />
            <Row align="middle" justify="center">
              <Col>
                <CameraFilled />
                <Select
                  style={{ width: 300, marginLeft: 10 }}
                  value={
                    settings.selectedWebcam || (settings.webcams.length ? settings.webcams[0].deviceId : undefined)
                  }
                  onSelect={(id) => {
                    settings.setSelectedWebcam(id as string);
                  }}
                  placeholder="Select a Camera"
                >
                  {(() => {
                    const opts: ReactElement[] = [];
                    for (const dev of settings.webcams) {
                      opts.push(
                        <Select.Option key={dev.deviceId} value={dev.deviceId}>
                          {dev.label}
                        </Select.Option>,
                      );
                    }
                    return opts;
                  })()}
                </Select>
              </Col>
            </Row>
            <br />
            <Row align="middle" justify="center">
              <Col>
                <AudioOutlined />
                <Select
                  style={{ width: 300, marginLeft: 10 }}
                  value={
                    settings.selectedMicrophone ||
                    (settings.microphones.length ? settings.microphones[0].deviceId : undefined)
                  }
                  onSelect={(id) => {
                    settings.setSelectedMicrophone(id as string);
                  }}
                  placeholder="Select a Microphone"
                >
                  {(() => {
                    const opts: ReactElement[] = [];
                    for (const dev of settings.microphones) {
                      opts.push(
                        <Select.Option key={dev.deviceId} value={dev.deviceId}>
                          {dev.label}
                        </Select.Option>,
                      );
                    }
                    return opts;
                  })()}
                </Select>
              </Col>
            </Row>
            <br />
          </Layout>
        </Modal>
        <StyledSider width={300}>
          {webcamDisplays.map((v, i) => {
            return (
              <StyledWebcam key={v.ref.key} index={i}>
                {v.ref}
                <div className="profile">
                  <UserAvatar profile={v.profile} />
                  <Typography.Text>
                    {v.profile.first_name} {v.profile.last_name}
                  </Typography.Text>
                </div>
              </StyledWebcam>
            );
          })}
          <Messaging messages={messages} setMessages={setMessages} height={webcamDisplays.length * webcamHeight} />
        </StyledSider>
        <Layout.Content>
          <StyledVideo
            autoPlay
            loop
            ref={(ref) => {
              screenRef.current = ref ?? undefined;
            }}
          />
        </Layout.Content>
        <StyledTools>
          <Tooltip title="Toggle Mute">
            <Button
              ghost={!micEnabled}
              onClick={() => setMicEnabled(!micEnabled)}
              size={'large'}
              shape="circle"
              style={{ margin: '0 10px' }}
            >
              <AudioOutlined size={20} style={{ color: micEnabled ? '#000' : '#fff' }} />
            </Button>
          </Tooltip>
          <Tooltip title="Toggle Webcam">
            <Button
              ghost={!webcamEnabled}
              onClick={() => setWebcamEnabled(!webcamEnabled)}
              size={'large'}
              shape="circle"
              style={{ margin: '0 10px' }}
            >
              <VideoCameraOutlined size={20} style={{ color: webcamEnabled ? '#000' : '#fff' }} />
            </Button>
          </Tooltip>
          <Tooltip title="Share Screen">
            <Button
              ghost={!screenEnabled}
              onClick={() => setScreenEnabled(!screenEnabled)}
              size={'large'}
              shape="circle"
              style={{ margin: '0 10px' }}
            >
              <DesktopOutlined size={20} style={{ color: screenEnabled ? '#000' : '#fff' }} />
            </Button>
          </Tooltip>
          <Tooltip title="Hang Up">
            <Button
              onClick={hangup}
              size={'large'}
              shape="circle"
              style={{ backgroundColor: '#c50505', margin: '0 10px' }}
            >
              <PhoneFilled size={20} rotate={225} style={{ color: '#fff' }} />
            </Button>
          </Tooltip>
        </StyledTools>
      </StyledLayout>
    </StyledLayout>
  );
}
