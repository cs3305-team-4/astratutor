
export interface Config {
    apiUrl: string;
}

const config: Config = {
    apiUrl: process.env.REACT_APP_API_URL as string,
}

console.log(process)

export default config;

console.log(config)