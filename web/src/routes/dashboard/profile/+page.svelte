<script>
	import { onMount } from 'svelte';
	import apiService from '$lib/api';

	// User data
	let user = $state();
	let isLoading = $state(true);
	let errorMessage = $state(null);

	onMount(async () => {
		await getProfile();
	});

	// Function untuk ambil data profile
	const getProfile = async () => {
		try {
			isLoading = true;
			const response = await apiService.get('/users/profile');
			if (!response.ok) {
				errorMessage = await response.text();
				isLoading = false;
				return;
			}
			user = await response.json();
		} catch (error) {
			errorMessage = 'Terjadi kesalahan saat mengambil data profil';
			isLoading = false;
			return;
		} finally {
			isLoading = false;
		}
	};

	function formatDate(isoString) {
		if (!isoString) return '-';
		const date = new Date(isoString);
		return new Intl.DateTimeFormat('id-ID', {
			dateStyle: 'long'
		}).format(date);
	}

	function getInitials(name = '') {
		if (!name) return '??';
		return name
			.split(' ')
			.map((n) => n[0])
			.join('')
			.toUpperCase()
			.slice(0, 2);
	}
</script>

<div class="mx-auto max-w-4xl">
	<!-- Page Header -->
	<div class="mb-8">
		<h1 class="text-3xl font-bold text-gray-800">Profil Pengguna</h1>
		<p class="text-gray-600">Informasi profil Anda</p>
	</div>

	<!-- Profile Card -->
	<div class="card border border-base-200 bg-base-100 shadow-xl">
		<div class="card-body">
			{#if isLoading}
				<!-- Skeleton UI while loading -->
				<div class="flex flex-col gap-8 md:flex-row">
					<!-- Avatar Skeleton -->
					<div class="flex flex-col items-center">
						<div class="avatar mb-4">
							<div class="h-24 w-24 animate-pulse rounded-full bg-base-300"></div>
						</div>
					</div>

					<!-- Profile Info Skeleton -->
					<div class="flex-1">
						<div class="grid grid-cols-1 gap-6 md:grid-cols-2">
							{#each Array(7) as _, i}
								<div>
									<div class="label">
										<div
											class="label-text h-4 w-24 animate-pulse rounded bg-base-300 font-medium"
										></div>
									</div>
									<div class="mt-1 h-6 w-full animate-pulse rounded bg-base-300"></div>
								</div>
							{/each}
						</div>
					</div>
				</div>
			{:else if errorMessage}
				<!-- Error Message -->
				<div class="alert alert-error">
					<div>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-6 w-6 shrink-0 stroke-current"
							fill="none"
							viewBox="0 0 24 24"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
							/>
						</svg>
						<span>{errorMessage}</span>
					</div>
				</div>
			{:else}
				<!-- Profile Content -->
				<div class="flex flex-col gap-8 md:flex-row">
					<!-- Avatar Section -->
					<div class="flex flex-col items-center">
						<div class="avatar mb-4">
							<div
								class="h-24 w-24 rounded-full ring ring-primary ring-offset-2 ring-offset-base-100"
							>
								<img
									src="https://placehold.co/400x400/ffffff/000000?text={getInitials(
										user.nama_lengkap
									)}"
									alt="User Avatar"
								/>
							</div>
						</div>
					</div>

					<!-- Profile Info Section -->
					<div class="flex-1">
						<div class="grid grid-cols-1 gap-6 md:grid-cols-2">
							<div>
								<label class="label" for="nama">
									<span class="label-text font-medium">Nama Lengkap</span>
								</label>
								<p class="text-lg" id="nama">{user?.nama_lengkap}</p>
							</div>

							<div>
								<label class="label" for="nip">
									<span class="label-text font-medium">NIP</span>
								</label>
								<p class="text-lg" id="nip">{user?.nip}</p>
							</div>

							<div>
								<label class="label" for="jabatan">
									<span class="label-text font-medium">Jabatan</span>
								</label>
								<p class="text-lg" id="jabatan">{user?.jabatan}</p>
							</div>

							<div>
								<label class="label" for="email">
									<span class="label-text font-medium">Email</span>
								</label>
								<p class="text-lg" id="email">{user?.email}</p>
							</div>

							<div>
								<label class="label" for="telepon">
									<span class="label-text font-medium">Telepon</span>
								</label>
								<p class="text-lg" id="telepon">{user?.nomor_telepon}</p>
							</div>

							<div>
								<label class="label" for="departemen">
									<span class="label-text font-medium">Unit Kerja</span>
								</label>
								<p id="departemen" class="text-lg">{user?.unit_kerja}</p>
							</div>

							<div>
								<label class="label" for="tanggal-bergabung">
									<span class="label-text font-medium">Tanggal Bergabung</span>
								</label>
								<p id="tanggal-bergabung" class="text-lg">{formatDate(user?.created_at)}</p>
							</div>
						</div>
					</div>
				</div>
			{/if}
		</div>
	</div>
</div>
