<script>
	import { goto } from '$app/navigation';
	import { onDestroy } from 'svelte';
	import apiService from '$lib/api';

	let countdown = $state(5);
	let timer;

	// Start countdown
	timer = setInterval(() => {
		countdown--;
		if (countdown <= 0) {
			clearInterval(timer);
			logoutNow();
		}
	}, 1000);

	// Cleanup timer on component destroy
	onDestroy(() => {
		if (timer) {
			clearInterval(timer);
		}
	});

	// Function to logout immediately
	const logoutNow = async () => {
		if (timer) {
			clearInterval(timer);
		}
		// Call logout API and redirect
		try {
			await apiService.post('/auth/logout');
		} catch (error) {
			console.error('Logout failed:', error);
		} finally {
			goto('/login');
		}
	};
</script>

<div
	class="flex min-h-screen items-center justify-center bg-gradient-to-b from-base-100 to-base-200 p-4"
>
	<div class="card w-full max-w-md border border-base-200 bg-base-100 shadow-xl">
		<div class="card-body items-center text-center">
			<div class="mb-4 rounded-full bg-primary/10 p-4">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="h-12 w-12 text-primary"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
					/>
				</svg>
			</div>

			<h2 class="card-title text-2xl font-bold">Berhasil Keluar</h2>
			<p class="mb-6 text-gray-600">
				Anda telah keluar dari sistem. Anda akan dialihkan ke halaman login.
			</p>

			<div class="mb-6">
				<p class="text-lg">Mengalihkan dalam <span class="font-bold">{countdown}</span> detik...</p>
			</div>

			<div class="card-actions w-full">
				<button class="btn w-full btn-primary" onclick={logoutNow}>
					<svg
						xmlns="http://www.w3.org/2000/svg"
						class="mr-2 h-5 w-5"
						fill="none"
						viewBox="0 0 24 24"
						stroke="currentColor"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
						/>
					</svg>
					Keluar Sekarang
				</button>

				<a href="/" class="btn w-full btn-link">Kembali ke Beranda</a>
			</div>
		</div>
	</div>
</div>
