import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';
import { loadEnv } from 'vite'

const env = loadEnv(process.env.NODE_ENV, process.cwd(), '')

// https://astro.build/config
export default defineConfig({
	trailingSlash: 'always',
	build: {
		format: 'directory'
	},
	site: 'https://email-package.tommymay.dev',
	integrations: [
		starlight({
			title: 'email',
			social: {
				github: 'https://github.com/itmayziii/email',
			},
			head: [
				{
					tag: 'script',
					attrs: {
						src: `https://www.googletagmanager.com/gtag/js?id=${env.GOOGLE_ANALYTICS_MEASUREMENT_ID}`,
						async: true
					},
				},
				{
					tag: 'script',
					content: `
  window.dataLayer = window.dataLayer || [];
  function gtag(){dataLayer.push(arguments);}
  gtag('js', new Date());

  gtag('config', '${env.GOOGLE_ANALYTICS_MEASUREMENT_ID}');					
					`
				}
			],
			sidebar: [
				{
					label: 'Start Here',
					items: [
						{ label: 'Quickstart', link: '/quickstart/' },
					],
				},
				{
					label: 'Guides',
					items: [
						{ label: 'Event Format', link: '/guides/event-format/' },
						{ label: 'Customize', link: '/guides/customize/' },
						{ label: 'Deploy', link: '/guides/deploy/' },
					],
				},
				{
					label: 'Reference',
					autogenerate: { directory: 'reference' },
				},
			],
		}),
	],
});
