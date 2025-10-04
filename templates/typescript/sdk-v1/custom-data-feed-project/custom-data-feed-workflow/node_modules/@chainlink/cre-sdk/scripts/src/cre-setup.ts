import { $ } from 'bun'

export const main = async () => {
	try {
		await $`bunx cre-setup`
	} catch {
		await $`bun --bun ../cre-sdk-javy-plugin/bin/setup.ts`
	}
}
