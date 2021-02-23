import {
  AudioOutlined,
  BgColorsOutlined,
  CameraFilled,
  DeleteOutlined,
  DesktopOutlined,
  EditOutlined,
  MinusOutlined,
  PhoneFilled,
  ScissorOutlined,
  SettingFilled,
  UndoOutlined,
  VideoCameraOutlined,
} from '@ant-design/icons';
import { Button, Col, Layout, Modal, Radio, Row, Select, Tooltip, Typography } from 'antd';
import React, { ReactElement, useContext, useEffect, useRef, useState } from 'react';
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
import { Stage, Layer, Image } from 'react-konva';
import { Stage as StageType, stages } from 'konva/types/Stage';
import Konva from 'konva';
import { Layer as LayerType } from 'konva/types/Layer';
import { SketchPicker } from 'react-color';
import { Collection } from 'konva/types/Util';
import { Line } from 'konva/types/shapes/Line';

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
  width: 100%;
  height: calc(100% - 88px);
`;

const StyledStreaming = styled.div`
  @keyframes transp {
    0% {
      opacity: 1;
    }
    50% {
      opacity: 1;
    }
    100% {
      opacity: 0;
    }
  }
  position: fixed;
  bottom: 120px;
  right: 40px;
  color: #fff;
  opacity: 0;
  animation: transp 8s;
`;

const StyledDrawMenu = styled.div`
  position: fixed;
  right: 30px;
  bottom: 24px;
`;

interface IJoin {
  layerJson: string;
  background: string;
}

function sleep(ms: number): unknown {
  return new Promise((resolve) => setTimeout(resolve, ms));
}
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
  const [streamingID, setStreamingID] = React.useState<string>('');
  const screenRef = useRef<HTMLVideoElement>();

  const [isPaint, setIsPaint] = useState(false);
  const [lastLine, setLastLine] = useState<Konva.Line>();
  const [mode, setMode] = useState<'brush' | 'line' | 'eraser'>('brush');
  const [bg, setBg] = useState('black');
  const stage = useRef<StageType>();
  const layer = useRef<LayerType>();

  const wipe = () => {
    layer.current?.removeChildren();
    layer.current?.clear();
  };

  const undo = () => {
    console.log(lastLine);
    if (layer.current) {
      layer.current.children.splice(layer.current.children.length - 1, 1);
      layer.current.batchDraw();
    }
  };

  const onDisconnect = (id: string) => {
    setWebcamDisplays((prev) => prev.filter((v) => v.profile.account_id !== id));
    setStreamingID((prev) => {
      console.log('Screen no longer receiving', prev);
      if (prev !== '') {
        prev = '';
        if (screenRef.current?.srcObject) {
          (screenRef.current.srcObject as MediaStream).getTracks().forEach((v) => {
            v.enabled = false;
            v.stop();
          });
          screenRef.current.srcObject = null;
        }
      }
      return prev;
    });
  };
  useAsync(async () => {
    // Signalling can be none if classroom page is refreshed before being sent back to lobby
    if (signalling == null) return;
    const credentials = await api.services.getTurnCredentials();
    handler.current = new WebRTCHandler(signalling, credentials, () => setAddingPeer(true));
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
        case MESSAGE_TYPE.STOP_STREAM: {
          if (message.dest !== signalling.id) return;
          console.log('MESSAGE STOP_STREAM');
          setScreenEnabled(false);
          break;
        }
        case MESSAGE_TYPE.DRAW: {
          if (layer.current && message.data) {
            console.log('MESSAGE DRAW', message.data);
            const line = new Konva.Line(JSON.parse(message.data));
            layer.current.add(line);
            layer.current.batchDraw();
          }
          break;
        }
        case MESSAGE_TYPE.UNDO: {
          undo();
          break;
        }
        case MESSAGE_TYPE.WIPE: {
          wipe();
          break;
        }
        case MESSAGE_TYPE.CHANGE_BG: {
          setBg(message.data);
          break;
        }
        case MESSAGE_TYPE.INIT: {
          const data = message.data as IJoin;
          if (lastLine === undefined) {
            console.log('INIT RECEIVED', message.data, stage.current);
            const children: string[] = JSON.parse(data.layerJson);
            const lines: Array<Line> = [];
            console.log(children);
            for (const child of Object.values(children)) {
              lines.push(new Konva.Line(JSON.parse(child)));
            }
            if (layer.current) {
              for (const child of lines) {
                layer.current.add(child);
              }
              layer.current.batchDraw();
            }
            setBg(data.background);
            setLastLine(new Konva.Line());
          }
          break;
        }
        case MESSAGE_TYPE.LEAVE: {
          onDisconnect(message.data);
          break;
        }
      }
    });
    signalling.send(MESSAGE_TYPE.AHOY_HOY, '', null);

    handler.current.ontrackremove = (id: string, correlation: StreamType, event: RTCTrackEvent) => {
      console.log(correlation);
      switch (correlation) {
        case StreamType.Screen:
          setStreamingID((prev) => {
            console.log('Screen no longer receiving', prev);
            if (prev !== '') {
              prev = '';
              if (screenRef.current?.srcObject) {
                (screenRef.current.srcObject as MediaStream).getTracks().forEach((v) => {
                  v.enabled = false;
                  v.stop();
                });
                screenRef.current.srcObject = null;
              }
            }
            return prev;
          });
          break;
      }
    };
    handler.current.ondisconnect = onDisconnect;
    handler.current.ontrack = (id: string, correlation: StreamType, event: RTCTrackEvent) => {
      console.log('NEW TRACK', id, correlation, event);
      const stream = event.streams.length ? event.streams[0] : new MediaStream();
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
        case StreamType.Screen:
          console.log('Other user started screensharing');
          if (screenRef.current) {
            ((screenRef.current.srcObject as MediaStream) || null)?.getTracks().forEach((v) => {
              v.enabled = false;
              v.stop();
              if (streamingID === '') {
                handler.current?.removeTrack(v);
              }
            });
            screenRef.current.srcObject = stream;
            screenRef.current.play();
          }
          setStreamingID(id);
          break;
      }
    };
  }, []);
  useEffect(() => console.log(webcamDisplays), [webcamDisplays]);

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
        muted
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
        if (layer.current) {
          const data: IJoin = {
            layerJson: JSON.stringify(layer.current.children),
            background: bg,
          };
          signalling?.send(MESSAGE_TYPE.INIT, '', data);
        }
        setAddingPeer(false);
      }
    }
  }, [settings.webcamStream, webcamEnabled, addingPeer, settings.selectedWebcam]);

  useAsync(async () => {
    if (screenEnabled) {
      const src = await screenStream();
      if (streamingID !== '') {
        console.log('stopping other stream');
        signalling?.send(MESSAGE_TYPE.STOP_STREAM, streamingID, { stop: true });
        await sleep(500);
      }
      if (!src) {
        setScreenEnabled(false);
        return;
      }
      setStreamingID('');
      src.onremovetrack = () => {
        setScreenEnabled(false);
      };
      for (const track of src.getTracks()) {
        track.onended = () => {
          setScreenEnabled(false);
        };
        handler.current?.addTrack(track, StreamType.Screen, src);
      }
      setScreen(src);
      if (screenRef.current) {
        screenRef.current.srcObject = src;
        screenRef.current.play();
      }
    } else {
      if (screenRef.current) {
        screenRef.current.srcObject = null;
      }
      screen?.getTracks().forEach((v) => {
        v.enabled = false;
        v.stop();
        handler.current?.removeTrack(v);
      });
      setScreen(undefined);
      setStreamingID('');
    }
  }, [screenEnabled]);

  const hangup = () => {
    handler.current?.close();
    settings.webcamStream?.getVideoTracks().forEach((v) => {
      v.stop();
    });
    signalling?.send(MESSAGE_TYPE.LEAVE, '', api.account?.id);
    history.push(`/lessons/${lid}/goodbye`);
  };

  useEffect(() => {
    settings.webcamStream?.getAudioTracks().forEach((v) => {
      v.enabled = micEnabled;
    });
  }, [micEnabled, settings.webcamStream]);

  useAsync(async () => {
    if (!isPaint) {
      console.log('SEND DRAW');
      signalling?.send(MESSAGE_TYPE.DRAW, '', lastLine?.toJSON());
    }
  }, [isPaint]);

  const [pick, setPick] = useState(false);
  const [color, setColor] = useState(
    '#' +
      Math.floor(Math.random() * 250 + 5).toString(16) +
      Math.floor(Math.random() * 250 + 5).toString(16) +
      Math.floor(Math.random() * 250 + 5).toString(16),
  );

  const [width, setWidth] = useState(window.innerWidth);
  const [height, setHeight] = useState(window.innerHeight);

  window.onresize = (e: UIEvent) => {
    setWidth(window.innerWidth);
    setHeight(window.innerHeight);
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
            style={{ backgroundColor: bg }}
            ref={(ref) => {
              screenRef.current = ref ?? undefined;
            }}
          />
          <Stage
            ref={(r) => {
              stage.current = r ?? undefined;
            }}
            onMouseDown={(e) => {
              console.log(e);
              if (!stage.current || !layer.current) return;
              setIsPaint(true);
              const pos = stage.current.getPointerPosition();
              if (!pos) return;
              const line = new Konva.Line({
                stroke: color,
                strokeWidth: 5,
                globalCompositeOperation: mode === 'brush' || mode === 'line' ? 'source-over' : 'destination-out',
                points: [pos.x, pos.y],
              });
              setLastLine(line);
              layer.current.add(line);
            }}
            width={width - 300}
            height={height - 90}
            style={{
              background: 'transparent',
              position: 'fixed',
              width: '100%',
              height: 'calc(100vh - 88px)',
              top: 0,
              cursor: 'crosshair',
            }}
            onMouseMove={(e) => {
              if (!isPaint) {
                return;
              }
              if (!stage.current || !layer.current || !lastLine) return;
              const pos = stage.current.getPointerPosition();
              if (!pos) return;
              let newPoints: number[];
              switch (mode) {
                case 'line':
                  newPoints = [lastLine.points()[0], lastLine.points()[1], pos.x, pos.y];
                  break;
                default:
                  newPoints = lastLine.points().concat([pos.x, pos.y]);
                  break;
              }
              lastLine.points(newPoints);
              layer.current.batchDraw();
            }}
            onMouseUp={(e) => {
              setIsPaint(false);
            }}
          >
            <Layer
              ref={(r) => {
                layer.current = r ?? undefined;
              }}
            ></Layer>
          </Stage>
          {streamingID.length > 0 && (
            <StyledStreaming>
              <UserAvatar
                profile={settings.otherProfiles[streamingID]}
                props={{ style: { marginRight: '.5em' }, size: 'small' }}
              />
              {settings.otherProfiles[streamingID].first_name} {settings.otherProfiles[streamingID].last_name}
            </StyledStreaming>
          )}
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
        <StyledDrawMenu>
          <Tooltip title="Free Draw">
            <Button
              ghost={mode !== 'brush'}
              onClick={() => setMode('brush')}
              size={'large'}
              style={{ margin: '0 3px' }}
            >
              <EditOutlined size={10} style={{ color: mode === 'brush' ? '#000' : '#fff' }} />
            </Button>
          </Tooltip>
          <Tooltip title="Line Draw">
            <Button ghost={mode !== 'line'} onClick={() => setMode('line')} size={'large'} style={{ margin: '0 3px' }}>
              <MinusOutlined size={10} style={{ color: mode === 'line' ? '#000' : '#fff' }} />
            </Button>
          </Tooltip>
          <Tooltip title="Eraser">
            <Button
              ghost={mode !== 'eraser'}
              onClick={() => setMode('eraser')}
              size={'large'}
              style={{ margin: '0 3px' }}
            >
              <ScissorOutlined size={10} style={{ color: mode === 'eraser' ? '#000' : '#fff' }} />
            </Button>
          </Tooltip>
          <Tooltip title="Wipe">
            <Button
              ghost={true}
              onClick={() => {
                wipe();
                signalling?.send(MESSAGE_TYPE.WIPE, '', {});
              }}
              size={'large'}
              style={{ margin: '0 3px' }}
            >
              <DeleteOutlined size={10} style={{ color: '#fff' }} />
            </Button>
          </Tooltip>
          <Tooltip title="Undo">
            <Button
              ghost={true}
              onClick={() => {
                undo();
                signalling?.send(MESSAGE_TYPE.UNDO, '', {});
              }}
              size={'large'}
              style={{ margin: '0 3px' }}
            >
              <UndoOutlined size={10} style={{ color: '#fff' }} />
            </Button>
          </Tooltip>
          <Tooltip title="Select colour">
            <Button ghost={!pick} onClick={() => setPick(!pick)} size={'large'} style={{ margin: '0 3px' }}>
              <BgColorsOutlined size={10} style={{ color: pick ? '#000' : '#fff' }} />
            </Button>
          </Tooltip>
          {pick && (
            <div style={{ position: 'fixed', bottom: 65, right: 0, zIndex: 10000 }}>
              <SketchPicker color={color} onChange={(e) => setColor(e.hex)} />
            </div>
          )}
          <Radio.Group
            value={bg}
            onChange={(e) => {
              setBg(e.target.value);
              signalling?.send(MESSAGE_TYPE.CHANGE_BG, '', e.target.value);
            }}
            style={{ position: 'fixed', left: 330, bottom: 27 }}
          >
            <Radio.Button style={{ background: 'transparent', color: '#fff' }} value="black">
              Black
            </Radio.Button>
            <Radio.Button style={{ background: 'transparent', color: '#fff' }} value="white">
              White
            </Radio.Button>
          </Radio.Group>
        </StyledDrawMenu>
      </StyledLayout>
    </StyledLayout>
  );
}
