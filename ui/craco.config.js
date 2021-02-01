const CracoLessPlugin = require('craco-less');

module.exports = {
  plugins: [
    {
      plugin: CracoLessPlugin,
      options: {
        lessLoaderOptions: {
          lessOptions: {
            modifyVars: { '@primary-color': 'rgba(255,215,0,0.75)' },
            javascriptEnabled: true,
          },
        },
      },
    },
  ],
};
