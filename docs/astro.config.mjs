import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

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
