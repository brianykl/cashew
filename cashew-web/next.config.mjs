/** @type {import('next').NextConfig} */
const nextConfig = {
    output: 'standalone',
    async headers() {
      return [
        {
          source: '/(.*)',
          headers: [
            {
              key: 'Content-Security-Policy',
              value: `
                default-src 'self' https://cdn.plaid.com;
                script-src 'self' 'unsafe-inline' 'unsafe-eval' https://cdn.plaid.com/link/v2/stable/link-initialize.js;
                frame-src 'self' https://cdn.plaid.com;
                connect-src 'self' https://production.plaid.com http://localhost:8080;
              `.replace(/\s{2,}/g, ' ').trim()
            }
          ]
        }
      ];
    }
  };
  
  export default nextConfig;
