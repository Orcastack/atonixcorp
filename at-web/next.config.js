/** @type {import('next').NextConfig} */
const nextConfig = {
  experimental: {
    turbo: {
      loaders: {
        css: false,
        postcss: false,
      },
    },
  },
};

module.exports = nextConfig;
