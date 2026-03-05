<script>
	// Mock data for statistics
	const stats = [
		{
			title: 'Total Situs',
			value: '52',
			icon: 'M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4',
			color: 'primary'
		},
		{
			title: 'Jenis Situs',
			value: '12',
			icon: 'M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z',
			color: 'secondary'
		},
		{
			title: 'Petugas Aktif',
			value: '8',
			icon: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z',
			color: 'accent'
		},
		{
			title: 'Laporan Bulan Ini',
			value: '24',
			icon: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z',
			color: 'info'
		}
	];

	// Mock recent activities
	const recentActivities = [
		{
			id: 1,
			user: 'Ahmad Fauzi',
			action: 'Menambahkan situs baru',
			target: 'Masjid Al-Muttaqin',
			time: '2 menit yang lalu',
			type: 'create'
		},
		{
			id: 2,
			user: 'Siti Nurhaliza',
			action: 'Memperbarui data',
			target: 'Gereja Katolik Santo Yosef',
			time: '15 menit yang lalu',
			type: 'update'
		},
		{
			id: 3,
			user: 'Budi Santoso',
			action: 'Menghapus situs',
			target: 'Pura Tirta Suci',
			time: '1 jam yang lalu',
			type: 'delete'
		},
		{
			id: 4,
			user: 'Lina Marlina',
			action: 'Melaporkan masalah',
			target: 'Vihara Dharma Bhakti',
			time: '2 jam yang lalu',
			type: 'report'
		},
		{
			id: 5,
			user: 'Rizki Pratama',
			action: 'Menyetujui situs',
			target: 'Klenteng Hok Tek Bio',
			time: '3 jam yang lalu',
			type: 'approve'
		}
	];

	// Mock sites data
	const sites = [
		{
			id: 1,
			name: 'Masjid Al-Muttaqin',
			type: 'Masjid',
			location: 'Desa Ciemas',
			status: 'Aktif',
			lastUpdated: '2023-05-15'
		},
		{
			id: 2,
			name: 'Gereja Katolik Santo Yosef',
			type: 'Gereja',
			location: 'Desa Citorek',
			status: 'Aktif',
			lastUpdated: '2023-05-14'
		},
		{
			id: 3,
			name: 'Pura Tirta Suci',
			type: 'Pura',
			location: 'Desa Cibodas',
			status: 'Non-Aktif',
			lastUpdated: '2023-05-10'
		},
		{
			id: 4,
			name: 'Vihara Dharma Bhakti',
			type: 'Vihara',
			location: 'Desa Cikahuripan',
			status: 'Aktif',
			lastUpdated: '2023-05-12'
		},
		{
			id: 5,
			name: 'Klenteng Hok Tek Bio',
			type: 'Klenteng',
			location: 'Desa Surade',
			status: 'Verifikasi',
			lastUpdated: '2023-05-13'
		}
	];
</script>

<div class="mx-auto max-w-7xl">
	<!-- Page Header -->
	<div class="mb-8">
		<h1 class="text-3xl font-bold text-gray-800">Dashboard Admin</h1>
		<p class="text-gray-600">
			Selamat datang di sistem pendataan situs keagamaan KUA Kecamatan Ciemas
		</p>
	</div>

	<!-- Stats Cards -->
	<div class="mb-8 grid grid-cols-1 gap-6 md:grid-cols-2 lg:grid-cols-4">
		{#each stats as stat (stat.title)}
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

	<!-- Charts and Recent Activity -->
	<div class="mb-8 grid grid-cols-1 gap-6 lg:grid-cols-3">
		<!-- Chart Placeholder -->
		<div class="card border border-base-200 bg-base-100 shadow-xl lg:col-span-2">
			<div class="card-body">
				<h2 class="card-title text-lg font-bold">Statistik Situs Berdasarkan Jenis</h2>
				<div class="flex h-64 items-center justify-center rounded-lg bg-base-200">
					<div class="text-center">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="mx-auto h-12 w-12 text-gray-400"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
							/>
						</svg>
						<p class="mt-2 text-gray-500">Grafik distribusi jenis situs keagamaan</p>
					</div>
				</div>
			</div>
		</div>

		<!-- Recent Activity -->
		<div class="card border border-base-200 bg-base-100 shadow-xl">
			<div class="card-body">
				<h2 class="card-title text-lg font-bold">Aktivitas Terbaru</h2>
				<div class="space-y-4">
					{#each recentActivities as activity (activity.id)}
						<div class="flex items-start space-x-3">
							<div class="avatar">
								<div
									class="flex h-8 w-8 items-center justify-center rounded-full bg-primary text-xs text-primary-content"
								>
									{activity.user.charAt(0)}
								</div>
							</div>
							<div class="min-w-0 flex-1">
								<p class="truncate text-sm">
									<span class="font-medium">{activity.user}</span>
									{activity.action} <span class="font-medium">{activity.target}</span>
								</p>
								<p class="text-xs text-gray-500">{activity.time}</p>
							</div>
						</div>
					{/each}
				</div>
			</div>
		</div>
	</div>

	<!-- Recent Sites Table -->
	<div class="card border border-base-200 bg-base-100 shadow-xl">
		<div class="card-body">
			<div class="mb-4 flex items-center justify-between">
				<h2 class="card-title text-lg font-bold">Daftar Situs Terbaru</h2>
				<a href="/dashboard/situs" class="btn btn-sm btn-primary">Lihat Semua</a>
			</div>
			<div class="overflow-x-auto">
				<table class="table table-zebra">
					<thead>
						<tr>
							<th>Nama Situs</th>
							<th>Jenis</th>
							<th>Lokasi</th>
							<th>Status</th>
							<th>Terakhir Diubah</th>
						</tr>
					</thead>
					<tbody>
						{#each sites as site (site.id)}
							<tr>
								<td>{site.name}</td>
								<td>
									<span class="badge badge-outline">{site.type}</span>
								</td>
								<td>{site.location}</td>
								<td>
									{#if site.status === 'Aktif'}
										<span class="badge badge-success">{site.status}</span>
									{:else if site.status === 'Non-Aktif'}
										<span class="badge badge-error">{site.status}</span>
									{:else}
										<span class="badge badge-warning">{site.status}</span>
									{/if}
								</td>
								<td>{site.lastUpdated}</td>
							</tr>
						{/each}
					</tbody>
				</table>
			</div>
		</div>
	</div>
</div>
