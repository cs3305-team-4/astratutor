import React from 'react';
import { Signalling } from '../webrtc/signalling';

export interface ISettings {
  signalling: Signalling | null;
  fullscreen: boolean;
  setFullscreen: (v: boolean) => void;
  webcams: MediaDeviceInfo[];
  setWebcams: (v: MediaDeviceInfo[]) => void;
  microphones: MediaDeviceInfo[];
  setMicrophones: (v: MediaDeviceInfo[]) => void;
  selectedWebcam: string;
  setSelectedWebcam: (v: string) => void;
  selectedMicrophone: string;
  setSelectedMicrophone: (v: string) => void;
  webcamStream: MediaStream | null;
  setWebcamStream: (v: MediaStream) => void;
}

export const SettingsCTX = React.createContext<ISettings>({
  signalling: null,
  fullscreen: false,
  setFullscreen: (v) => null,
  webcams: [],
  setWebcams: (v) => null,
  microphones: [],
  setMicrophones: (v) => null,
  selectedWebcam: '',
  setSelectedWebcam: (v) => null,
  selectedMicrophone: '',
  setSelectedMicrophone: (v) => null,
  webcamStream: null,
  setWebcamStream: (v) => null,
});
