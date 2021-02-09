import React from 'react';

export interface ISettings {
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
});
