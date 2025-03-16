import { redirect } from '@sveltejs/kit';

export function load() {
	// The root URL should show the landing page content
	// This is a no-op redirect since the (landing) route group already handles the root URL
	return {};
} 