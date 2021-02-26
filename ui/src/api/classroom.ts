import React from 'react';
import { Signalling } from '../webrtc/signalling';
import { ProfileResponseDTO } from './definitions';

export interface ISettings {
  signalling?: Signalling;
  fullscreen: boolean;
  setFullscreen: (v: boolean) => void;
  webcams: MediaDeviceInfo[];
  setWebcams: React.Dispatch<React.SetStateAction<MediaDeviceInfo[]>>;
  microphones: MediaDeviceInfo[];
  setMicrophones: (v: MediaDeviceInfo[]) => void;
  selectedWebcam: string;
  setSelectedWebcam: (v: string) => void;
  selectedMicrophone: string;
  setSelectedMicrophone: (v: string) => void;
  webcamStream: MediaStream | null;
  setWebcamStream: (v: MediaStream) => void;
  otherProfiles: { [id: string]: ProfileResponseDTO };
}

export const SettingsCTX = React.createContext<ISettings>({
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
  otherProfiles: {},
});
