<script>
	import { page } from '$app/stores';

	// Mengambil nilai status dan pesan error dari SvelteKit page store
	let status = $derived($page.status);
	let errorMessage = $derived($page.error?.message || 'Terjadi kesalahan yang tidak diketahui');

	function getTitle() {
		switch (status) {
			case 404:
				return 'Halaman Tidak Ditemukan';
			case 500:
				return 'Kesalahan Server';
			case 403:
				return 'Akses Ditolak';
			case 401:
				return 'Belum Login';
			default:
				return 'Terjadi Kesalahan';
		}
	}
</script>

<div class="flex min-h-screen items-center justify-center bg-base-200">
	<div class="px-4 text-center">
		<h1 class="text-9xl font-bold text-primary">{status}</h1>
		<h2 class="mt-4 mb-2 text-3xl font-semibold">{getTitle()}</h2>
		<p class="mb-8 text-gray-600">{errorMessage}</p>
		<div class="flex flex-wrap justify-center gap-4">
			<a href="/dashboard" class="btn btn-primary">
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="mr-1 h-5 w-5"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6"
					/>
				</svg>
				Kembali ke Dashboard
			</a>
			<button class="btn btn-outline" onclick={() => window.history.back()}>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					class="mr-1 h-5 w-5"
					fill="none"
					viewBox="0 0 24 24"
					stroke="currentColor"
				>
					<path
						stroke-linecap="round"
						stroke-linejoin="round"
						stroke-width="2"
						d="M10 19l-7-7m0 0l7-7m-7 7h18"
					/>
				</svg>
				Kembali
			</button>
		</div>
	</div>
</div>
