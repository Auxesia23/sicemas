<script>
	import { onMount } from 'svelte';
	import apiService from '$lib/api'; // Sesuaikan path ini jika beda

	// State untuk menampung data dari API
	let loading = $state(true);
	let error = $state(null);

	let dashboardStats = $state({
		total_situs: 0,
		total_jenis: 0,
		petugas_aktif: 0,
		menunggu_verifikasi: 0
	});

	let chartData = $state([]); // Disiapkan buat dipakai Chart nanti
	let recentActivities = $state([]);
	let recentSites = $state([]);

	// Derived state untuk menggabungkan angka dari API dengan desain UI kartu
	let statsConfig = $derived([
		{
			title: 'Total Situs',
			value: dashboardStats.total_situs,
			icon: 'M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4',
			color: 'primary'
		},
		{
			title: 'Jenis Situs',
			value: dashboardStats.total_jenis,
			icon: 'M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z',
			color: 'secondary'
		},
		{
			title: 'Petugas Aktif',
			value: dashboardStats.petugas_aktif,
			icon: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z',
			color: 'accent'
		},
		{
			title: 'Menunggu Verifikasi',
			value: dashboardStats.menunggu_verifikasi,
			icon: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z',
			color: 'warning' // Diganti warning (kuning) biar lebih mencolok
		}
	]);

	// Fetch API on Mount
	onMount(async () => {
		try {
			loading = true;
			// Sesuaikan endpoint ini dengan route API Go lu
			const response = await apiService.get('/dashboard');

			if (!response.ok) {
				throw new Error('Gagal memuat data dashboard');
			}

			const data = await response.json();

			dashboardStats = data.stats;
			chartData = data.statistik_jenis;
			recentActivities = data.recent_activities || [];
			recentSites = data.recent_sites || [];
		} catch (err) {
			console.error(err);
			error = err.message;
		} finally {
			loading = false;
		}
	});

	// Helper: Mengubah timestamp jadi teks "x menit yang lalu"
	function timeAgo(dateString) {
		if (!dateString) return '-';
		const date = new Date(dateString);
		const now = new Date();
		const seconds = Math.round((now - date) / 1000);
		const minutes = Math.round(seconds / 60);
		const hours = Math.round(minutes / 60);
		const days = Math.round(hours / 24);

		if (seconds < 60) return 'Baru saja';
		if (minutes < 60) return `${minutes} menit yang lalu`;
		if (hours < 24) return `${hours} jam yang lalu`;
		if (days === 1) return 'Kemarin';
		return `${days} hari yang lalu`;
	}

	// Helper: Format tanggal untuk tabel
	function formatDate(dateString) {
		if (!dateString) return '-';
		const date = new Date(dateString);
		return date.toLocaleDateString('id-ID', {
			day: 'numeric',
			month: 'short',
			year: 'numeric'
		});
	}

	// Helper: Class warna badge sesuai status DB
	function getStatusBadgeClass(status) {
		switch (status?.toLowerCase()) {
			case 'menunggu':
				return 'badge-warning';
			case 'terverifikasi':
				return 'badge-success';
			case 'ditolak':
				return 'badge-error';
			default:
				return 'badge-ghost';
		}
	}
</script>

<div class="mx-auto max-w-7xl">
	<div class="mb-8">
		<h1 class="text-3xl font-bold text-gray-800">Dashboard Admin</h1>
		<p class="text-gray-600">
			Selamat datang di sistem pendataan situs keagamaan KUA Kecamatan Ciemas
		</p>
	</div>

	{#if loading}
		<div class="flex h-64 flex-col items-center justify-center">
			<span class="loading loading-lg loading-spinner text-primary"></span>
			<p class="mt-4 text-gray-500">Memuat data dashboard...</p>
		</div>
	{:else if error}
		<div class="mb-8 alert alert-error shadow-sm">
			<span>{error}</span>
			<button class="btn btn-sm" onclick={() => location.reload()}>Coba Lagi</button>
		</div>
	{:else}
		<div class="mb-8 grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-4">
			{#each statsConfig as stat (stat.title)}
				<div class="card border border-base-200 bg-base-100 shadow-xl">
					<div class="card-body">
						<div class="flex items-center justify-between">
							<div>
								<p class="text-sm text-gray-500">{stat.title}</p>
								<h3 class="text-3xl font-bold text-{stat.color}">{stat.value}</h3>
							</div>
							<div class={`bg-${stat.color}-100 rounded-full p-3`}>
								<svg
									xmlns="http://www.w3.org/2000/svg"
									class={`h-6 w-6 text-${stat.color}`}
									fill="none"
									viewBox="0 0 24 24"
									stroke="currentColor"
								>
									<path
										stroke-linecap="round"
										stroke-linejoin="round"
										stroke-width="2"
										d={stat.icon}
									/>
								</svg>
							</div>
						</div>
					</div>
				</div>
			{/each}
		</div>

		<div class="mb-8 grid grid-cols-1 gap-6 lg:grid-cols-3">
			<div class="card border border-base-200 bg-base-100 shadow-xl lg:col-span-2">
				<div class="card-body">
					<h2 class="card-title text-lg font-bold">Statistik Situs Berdasarkan Jenis</h2>

					<div class="mt-4 flex flex-col justify-center space-y-5 py-2">
						{#if chartData.length > 0}
							{@const maxSitus = Math.max(...chartData.map((d) => d.jumlah_situs))}

							{#each chartData as item}
								<div class="flex flex-col gap-1.5">
									<div class="flex justify-between text-sm">
										<span class="font-medium text-gray-700">{item.nama}</span>
										<span class="font-bold text-gray-900">{item.jumlah_situs} Situs</span>
									</div>
									<div class="h-3 w-full overflow-hidden rounded-full bg-base-200">
										<div
											class="h-full rounded-full bg-primary transition-all duration-1000 ease-out"
											style="width: {maxSitus === 0 ? 0 : (item.jumlah_situs / maxSitus) * 100}%"
										></div>
									</div>
								</div>
							{/each}
						{:else}
							<div class="py-10 text-center text-gray-400">
								Belum ada data statistik jenis situs.
							</div>
						{/if}
					</div>
				</div>
			</div>

			<div class="card border border-base-200 bg-base-100 shadow-xl">
				<div class="card-body">
					<h2 class="card-title text-lg font-bold">Aktivitas Terbaru</h2>
					<div class="space-y-4">
						{#each recentActivities as activity (activity.id)}
							<div class="flex items-start space-x-3">
								<div class="avatar">
									<div
										class="flex h-8 w-8 items-center justify-center rounded-full bg-primary text-xs font-bold text-primary-content"
									>
										{activity.user.charAt(0).toUpperCase()}
									</div>
								</div>
								<div class="min-w-0 flex-1">
									<p class="text-sm">
										<span class="font-medium text-gray-800">{activity.user}</span>
										<span class="text-gray-600">{activity.action}</span>
										<span class="font-medium text-gray-800">{activity.target}</span>
									</p>
									<p class="text-xs text-gray-500">{timeAgo(activity.created_at)}</p>
								</div>
							</div>
						{:else}
							<div class="text-center py-4 text-gray-500 text-sm">Belum ada aktivitas terbaru.</div>
						{/each}
					</div>
				</div>
			</div>
		</div>

		<div class="card border border-base-200 bg-base-100 shadow-xl">
			<div class="card-body">
				<div class="mb-4 flex items-center justify-between">
					<h2 class="card-title text-lg font-bold">Daftar Situs Terbaru</h2>
					<a href="/dashboard/situs" class="btn btn-sm btn-primary">Lihat Semua</a>
				</div>
				<div class="overflow-x-auto">
					<table class="table w-full table-zebra">
						<thead>
							<tr>
								<th>Nama Situs</th>
								<th>Jenis</th>
								<th>Lokasi</th>
								<th>Pendata</th>
								<th>Status</th>
								<th>Terakhir Diubah</th>
							</tr>
						</thead>
						<tbody>
							{#each recentSites as site (site.id)}
								<tr>
									<td class="font-medium">{site.nama}</td>
									<td>
										<span class="badge badge-outline badge-sm">{site.jenis}</span>
									</td>
									<td>{site.lokasi}</td>
									<td>{site.pendata}</td>
									<td>
										<span class="badge badge-sm {getStatusBadgeClass(site.status)}">
											{site.status}
										</span>
									</td>
									<td class="text-sm text-gray-600">{formatDate(site.updated_at)}</td>
								</tr>
							{:else}
								<tr>
									<td colspan="6" class="text-center py-6 text-gray-500">
										Belum ada data situs yang terdaftar.
									</td>
								</tr>
							{/each}
						</tbody>
					</table>
				</div>
			</div>
		</div>
	{/if}
</div>
