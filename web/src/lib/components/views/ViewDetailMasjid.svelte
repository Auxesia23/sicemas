<script>
	// Svelte 5 Runes: Receive detail prop from parent
	let { detail } = $props();
</script>

<!-- Detail Masjid View -->
<div class="space-y-6">
	<!-- Perpustakaan -->
	{#if detail.perpustakaan}
		<div class="card border border-base-200 bg-base-100 shadow-lg">
			<div class="card-body p-4 sm:p-6">
				<h2 class="mb-4 card-title text-lg">Perpustakaan</h2>
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
					<div>
						<p class="mb-1 block text-xs font-medium text-base-content/70">Luas</p>
						<p class="font-medium">
							{detail.perpustakaan.luas_m2 ? `${detail.perpustakaan.luas_m2} m²` : '-'}
						</p>
					</div>
					<div>
						<p class="mb-1 block text-xs font-medium text-base-content/70">Jumlah Pengurus</p>
						<p class="font-medium">{detail.perpustakaan.jumlah_pengurus || '-'}</p>
					</div>
					<div>
						<p class="mb-1 block text-xs font-medium text-base-content/70">Kondisi</p>
						<p class="font-medium">
							<span
								class="badge {detail.perpustakaan.kondisi === 'Baik'
									? 'badge-success'
									: 'badge-ghost'}"
							>
								{detail.perpustakaan.kondisi || '-'}
							</span>
						</p>
					</div>
				</div>

				{#if detail.perpustakaan.jenis_buku && detail.perpustakaan.jenis_buku.length > 0}
					<div class="mt-4">
						<p class="mb-1 block text-xs font-medium text-base-content/70">Jenis Buku</p>
						<div class="mt-2 flex flex-wrap gap-2">
							{#each detail.perpustakaan.jenis_buku as buku (buku)}
								<span class="badge badge-primary">{buku}</span>
							{/each}
						</div>
					</div>
				{/if}
			</div>
		</div>
	{/if}

	<!-- Kalibrasi Arah Kiblat -->
	{#if detail.kalibrasi_arah_kiblat}
		<div class="card border border-base-200 bg-base-100 shadow-lg">
			<div class="card-body p-4 sm:p-6">
				<h2 class="mb-4 card-title text-lg">Kalibrasi Arah Kiblat</h2>
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
					<div>
						<p class="mb-1 block text-xs font-medium text-base-content/70">Azimut</p>
						<p class="font-mono font-medium">
							{detail.kalibrasi_arah_kiblat.azimut || '-'}
						</p>
					</div>
					<div>
						<p class="mb-1 block text-xs font-medium text-base-content/70">Tanggal Kalibrasi</p>
						<p class="font-medium">
							{detail.kalibrasi_arah_kiblat.tanggal_kalibrasi || '-'}
						</p>
					</div>
				</div>
			</div>
		</div>
	{/if}

	<!-- SDM Masjid -->
	{#if detail.sdm_masjid}
		<div class="card border border-base-200 bg-base-100 shadow-lg">
			<div class="card-body p-4 sm:p-6">
				<h2 class="mb-4 card-title text-lg">SDM Masjid</h2>
				<div class="grid grid-cols-2 gap-4 sm:grid-cols-3 lg:grid-cols-6">
					<div class="stat rounded-xl bg-base-200 p-3">
						<div class="stat-title text-xs">Imam</div>
						<div class="stat-value text-primary">{detail.sdm_masjid.jumlah_imam || 0}</div>
					</div>
					<div class="stat rounded-xl bg-base-200 p-3">
						<div class="stat-title text-xs">Khotib</div>
						<div class="stat-value text-secondary">{detail.sdm_masjid.jumlah_khotib || 0}</div>
					</div>
					<div class="stat rounded-xl bg-base-200 p-3">
						<div class="stat-title text-xs">Muadzin</div>
						<div class="stat-value text-accent">{detail.sdm_masjid.jumlah_muadzin || 0}</div>
					</div>
					<div class="stat rounded-xl bg-base-200 p-3">
						<div class="stat-title text-xs">Pengurus</div>
						<div class="stat-value text-info">{detail.sdm_masjid.jumlah_pengurus || 0}</div>
					</div>
					<div class="stat rounded-xl bg-base-200 p-3">
						<div class="stat-title text-xs">Remaja</div>
						<div class="stat-value text-warning">{detail.sdm_masjid.jumlah_remaja || 0}</div>
					</div>
					<div class="stat rounded-xl bg-base-200 p-3">
						<div class="stat-title text-xs">Jemaah</div>
						<div class="stat-value text-success">{detail.sdm_masjid.jumlah_jemaah || 0}</div>
					</div>
				</div>
			</div>
		</div>
	{/if}

	<!-- Pengurus DKM -->
	{#if detail.pengurus_dkm}
		<div class="card border border-base-200 bg-base-100 shadow-lg">
			<div class="card-body p-4 sm:p-6">
				<h2 class="mb-4 card-title text-lg">Pengurus DKM</h2>
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
					<div>
						<p class="mb-1 block text-xs font-medium text-base-content/70">Ketua</p>
						<p class="font-medium">{detail.pengurus_dkm.ketua || '-'}</p>
					</div>
					<div>
						<p class="mb-1 block text-xs font-medium text-base-content/70">Sekretaris</p>
						<div class="mt-2 flex flex-wrap gap-2">
							{#each detail.pengurus_dkm.sekretaris as sekretaris (sekretaris)}
								<span class="badge badge-outline">{sekretaris}</span>
							{/each}
						</div>
					</div>
					<div>
						<p class="mb-1 block text-xs font-medium text-base-content/70">Bendahara</p>
						<div class="mt-2 flex flex-wrap gap-2">
							{#each detail.pengurus_dkm.bendahara as bendahara (bendahara)}
								<span class="badge badge-outline">{bendahara}</span>
							{/each}
						</div>
					</div>
				</div>

				{#if detail.pengurus_dkm.nama_imam && detail.pengurus_dkm.nama_imam.length > 0}
					<div class="mt-4">
						<p class="mb-1 block text-xs font-medium text-base-content/70">Nama Imam</p>
						<div class="mt-2 flex flex-wrap gap-2">
							{#each detail.pengurus_dkm.nama_imam as imam (imam)}
								<span class="badge badge-outline">{imam}</span>
							{/each}
						</div>
					</div>
				{/if}

				{#if detail.pengurus_dkm.nama_muazdin && detail.pengurus_dkm.nama_muazdin.length > 0}
					<div class="mt-4">
						<p class="mb-1 block text-xs font-medium text-base-content/70">Nama Muadzin</p>
						<div class="mt-2 flex flex-wrap gap-2">
							{#each detail.pengurus_dkm.nama_muazdin as muazdin (muazdin)}
								<span class="badge badge-outline">{muazdin}</span>
							{/each}
						</div>
					</div>
				{/if}
			</div>
		</div>
	{/if}

	<!-- Fasilitas -->
	{#if detail.fasilitas_umum || detail.fasilitas_raman_anak || detail.fasilitas_disabilitas}
		<div class="card border border-base-200 bg-base-100 shadow-lg">
			<div class="card-body p-4 sm:p-6">
				<h2 class="mb-4 card-title text-lg">Fasilitas</h2>
				<div class="space-y-4">
					{#if detail.fasilitas_umum && detail.fasilitas_umum.length > 0}
						<div>
							<p class="mb-1 block text-xs font-medium text-base-content/70">Fasilitas Umum</p>
							<div class="mt-2 flex flex-wrap gap-2">
								{#each detail.fasilitas_umum as fasilitas (fasilitas)}
									<span class="badge badge-secondary">{fasilitas}</span>
								{/each}
							</div>
						</div>
					{/if}

					{#if detail.fasilitas_raman_anak && detail.fasilitas_raman_anak.length > 0}
						<div>
							<p class="mb-1 block text-xs font-medium text-base-content/70">
								Fasilitas Ramah Anak
							</p>
							<div class="mt-2 flex flex-wrap gap-2">
								{#each detail.fasilitas_raman_anak as fasilitas (fasilitas)}
									<span class="badge badge-accent">{fasilitas}</span>
								{/each}
							</div>
						</div>
					{/if}

					{#if detail.fasilitas_disabilitas && detail.fasilitas_disabilitas.length > 0}
						<div>
							<p class="mb-1 block text-xs font-medium text-base-content/70">
								Fasilitas Disabilitas
							</p>
							<div class="mt-2 flex flex-wrap gap-2">
								{#each detail.fasilitas_disabilitas as fasilitas (fasilitas)}
									<span class="badge badge-outline">{fasilitas}</span>
								{/each}
							</div>
						</div>
					{:else}
						<p class="text-sm text-base-content/50">Tidak ada fasilitas disabilitas</p>
					{/if}
				</div>
			</div>
		</div>
	{/if}

	<!-- Kegiatan -->
	{#if detail.kegiatan && detail.kegiatan.length > 0}
		<div class="card border border-base-200 bg-base-100 shadow-lg">
			<div class="card-body p-4 sm:p-6">
				<h2 class="mb-4 card-title text-lg">Kegiatan</h2>
				<div class="flex flex-wrap gap-2">
					{#each detail.kegiatan as kegiatan (kegiatan)}
						<span class="badge badge-success">{kegiatan}</span>
					{/each}
				</div>
			</div>
		</div>
	{/if}
</div>
