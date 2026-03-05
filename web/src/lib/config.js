// API Configuration
const config = {
	// Base URL for API calls - can be overridden by environment variables
	apiUrl: import.meta.env.VITE_API_URL,

	// Timeout for API requests (in milliseconds)
	timeout: 10000,

	// Headers that should be sent with every request
	defaultHeaders: {
		'Content-Type': 'application/json'
	},

	// Default fetch options
	defaultFetchOptions: {
		credentials: 'include' // Include cookies in requests
	}
};

export default config;
