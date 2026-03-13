<script>
	// Svelte 5 Runes: Receive detail prop from parent
	let { detail } = $props();
</script>

<!-- Detail Majelis Taklim View -->
<div class="space-y-6">
	<!-- Pengurus -->
	{#if detail.pengurus}
		<div class="card border border-base-200 bg-base-100 shadow-lg">
			<div class="card-body p-4 sm:p-6">
				<h2 class="mb-4 card-title text-lg">Pengurus</h2>
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
					<div>
						<p class="mb-1 block text-xs font-medium text-base-content/70">Ketua</p>
						<p class="font-medium">{detail.pengurus.ketua || '-'}</p>
					</div>
					<div>
						<p class="mb-1 block text-xs font-medium text-base-content/70">Sekretaris</p>
						<p class="font-medium">{detail.pengurus.sekretaris || '-'}</p>
					</div>
					<div>
						<p class="mb-1 block text-xs font-medium text-base-content/70">Bendahara</p>
						<p class="font-medium">{detail.pengurus.bendahara || '-'}</p>
					</div>
				</div>
			</div>
		</div>
	{/if}

	<!-- Kehadiran Pria -->
	{#if detail.kehadiran_pria}
		<div class="card border border-base-200 bg-base-100 shadow-lg">
			<div class="card-body p-4 sm:p-6">
				<h2 class="mb-4 card-title text-lg">Kehadiran Pria</h2>
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
					<div class="stat rounded-xl bg-base-200 p-3">
						<div class="stat-title text-xs">Dewasa</div>
						<div class="stat-value text-primary">{detail.kehadiran_pria.jumlah_dewasa ?? 0}</div>
					</div>
					<div class="stat rounded-xl bg-base-200 p-3">
						<div class="stat-title text-xs">Remaja</div>
						<div class="stat-value text-secondary">{detail.kehadiran_pria.jumlah_remaja ?? 0}</div>
					</div>
					<div class="flex flex-col justify-center">
						<p class="mb-1 block text-xs font-medium text-base-content/70">Waktu Pengajian</p>
						<p class="mt-1 font-medium">
							<span class="badge badge-primary">
								{detail.kehadiran_pria.waktu_pengajian || '-'}
							</span>
						</p>
					</div>
				</div>
			</div>
		</div>
	{/if}

	<!-- Kehadiran Wanita -->
	{#if detail.kehadiran_wanita}
		<div class="card border border-base-200 bg-base-100 shadow-lg">
			<div class="card-body p-4 sm:p-6">
				<h2 class="mb-4 card-title text-lg">Kehadiran Wanita</h2>
				<div class="grid grid-cols-1 gap-4 sm:grid-cols-3">
					<div class="stat rounded-xl bg-base-200 p-3">
						<div class="stat-title text-xs">Dewasa</div>
						<div class="stat-value text-accent">{detail.kehadiran_wanita.jumlah_dewasa ?? 0}</div>
					</div>
					<div class="stat rounded-xl bg-base-200 p-3">
						<div class="stat-title text-xs">Remaja</div>
						<div class="stat-value text-info">{detail.kehadiran_wanita.jumlah_remaja ?? 0}</div>
					</div>
					<div class="flex flex-col justify-center">
						<p class="mb-1 block text-xs font-medium text-base-content/70">Waktu Pengajian</p>
						<p class="mt-1 font-medium">
							<span class="badge badge-accent">
								{detail.kehadiran_wanita.waktu_pengajian || '-'}
							</span>
						</p>
					</div>
				</div>
			</div>
		</div>
	{/if}

	<!-- Daftar Penceramah -->
	{#if detail.penceramah && detail.penceramah.length > 0}
		<div class="card border border-base-200 bg-base-100 shadow-lg">
			<div class="card-body p-4 sm:p-6">
				<h2 class="mb-4 card-title text-lg">Daftar Penceramah</h2>
				<div class="space-y-3">
					{#each detail.penceramah as penceramah (penceramah.nama)}
						<div
							class="flex items-center justify-between rounded-lg border border-base-300 bg-base-200 p-4"
						>
							<div class="flex-1">
								<div class="text-base font-semibold">{penceramah.nama}</div>
								<div class="text-sm text-base-content/70">{penceramah.materi}</div>
							</div>
						</div>
					{/each}
				</div>
			</div>
		</div>
	{/if}
</div>
