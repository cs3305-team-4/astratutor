export interface Config {
  apiUrl: string;
  stripePublishableKey: string;
}

const config: Config = {
  apiUrl: process.env.REACT_APP_API_URL as string,
  stripePublishableKey: process.env.REACT_APP_STRIPE_PUBLISHABLE_KEY as string,
};

export default config;
