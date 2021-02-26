export interface Config {
  apiUrl: string;
  signallingUrl: string;
}

const config: Config = {
  apiUrl: process.env.REACT_APP_API_URL as string,
  signallingUrl: process.env.REACT_APP_SIGNALLING_URL as string,
};

export default config;
