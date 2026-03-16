import { PUBLIC_API_URL } from '$env/static/public';
// API Configuration
const config = {
	// Base URL for API calls - can be overridden by environment variables
	apiUrl: PUBLIC_API_URL,

	// Timeout for API requests (in milliseconds) - default 30 seconds
	timeout: 30000,

	// Headers that should be sent with every request
	// Note: Content-Type is handled dynamically based on request body
	defaultHeaders: {},

	// Default fetch options
	defaultFetchOptions: {
		credentials: 'include' // Include cookies in requests
	}
};

export default config;
