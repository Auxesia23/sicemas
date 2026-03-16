<script>
	import { onMount } from 'svelte';
	import { hasAnyPermission } from '$lib/permissions';

	let { data, children } = $props();
	let user = $derived(data.user);

	let isSidebarOpen = $state(false);

	// Toggle sidebar on mobile
	const toggleSidebar = () => {
		isSidebarOpen = !isSidebarOpen;
	};

	// Close sidebar when clicking outside on mobile
	const closeSidebar = (event) => {
		if (window.innerWidth < 768 && isSidebarOpen) {
			const sidebar = document.getElementById('sidebar');
			const menuButton = document.getElementById('menu-button');

			if (sidebar && menuButton) {
				if (!sidebar.contains(event.target) && !menuButton.contains(event.target)) {
					isSidebarOpen = false;
				}
			}
		}
	};

	onMount(() => {
		document.addEventListener('click', closeSidebar);

		// Handle window resize to close sidebar on larger screens
		const handleResize = () => {
			if (window.innerWidth >= 768) {
				isSidebarOpen = false;
			}
		};

		window.addEventListener('resize', handleResize);

		return () => {
			document.removeEventListener('click', closeSidebar);
			window.removeEventListener('resize', handleResize);
		};
	});

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

<div class="drawer lg:drawer-open">
	<input id="sidebar-drawer" type="checkbox" class="drawer-toggle" bind:checked={isSidebarOpen} />
	<div class="drawer-content flex flex-col">
		<!-- Top Navigation Bar -->
		<header class="navbar sticky top-0 z-50 border-b border-base-200 bg-base-100 px-4 py-3">
			<div class="flex w-full items-center justify-between">
				<div class="flex items-center">
					<label
						for="sidebar-drawer"
						class="drawer-button btn mr-2 btn-square btn-ghost lg:hidden"
						id="menu-button"
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-6 w-6"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M4 6h16M4 12h16M4 18h16"
							/>
						</svg>
					</label>
					<a href="/dashboard" class="btn text-xl font-bold text-primary btn-ghost">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="mr-2 h-6 w-6"
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
						Dashboard
					</a>
				</div>

				<div class="flex items-center space-x-4">
					<!-- User Profile Dropdown -->
					<div class="dropdown dropdown-end">
						<div tabindex="0" role="button" class="btn avatar btn-circle btn-ghost">
							<div class="w-10 rounded-full">
								<img
									src="https://placehold.co/400x400/ffffff/000000?text={getInitials(
										user.nama_lengkap
									)}"
									alt="User Avatar"
								/>
							</div>
						</div>
						<ul
							class="dropdown-content menu z-1 mt-3 w-52 menu-sm rounded-box bg-base-100 p-2 shadow"
						>
							<li>
								<a href="/dashboard/profile"
									><svg
										xmlns="http://www.w3.org/2000/svg"
										class="h-5 w-5"
										fill="none"
										viewBox="0 0 24 24"
										stroke="currentColor"
										><path
											stroke-linecap="round"
											stroke-linejoin="round"
											stroke-width="2"
											d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
										/></svg
									> Profile</a
								>
							</li>
							<li>
								<a href="/logout"
									><svg
										xmlns="http://www.w3.org/2000/svg"
										class="h-5 w-5"
										fill="none"
										viewBox="0 0 24 24"
										stroke="currentColor"
										><path
											stroke-linecap="round"
											stroke-linejoin="round"
											stroke-width="2"
											d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
										/></svg
									> Logout</a
								>
							</li>
						</ul>
					</div>
				</div>
			</div>
		</header>

		<!-- Main Content -->
		<main class="min-h-[calc(100vh-4rem)] flex-1 bg-base-200 p-4 md:p-6">
			{@render children()}
		</main>
	</div>

	<!-- Sidebar Navigation -->
	<div id="sidebar" class="drawer-side z-40">
		<label for="sidebar-drawer" class="drawer-overlay"></label>
		<aside
			class="menu min-h-screen w-64 border-r border-base-200 bg-base-100 p-0 text-base-content"
		>
			<div class="border-b border-base-200 p-4">
				<div class="flex items-center space-x-3 p-2">
					<div class="avatar">
						<div class="w-10 rounded-full">
							<img src="/kemenag.webp" alt="KUA Logo" />
						</div>
					</div>
					<div>
						<h3 class="text-lg font-bold">KUA Ciemas</h3>
						<p class="text-xs opacity-70">Admin Dashboard</p>
					</div>
				</div>
			</div>

			<ul class="pt-2">
				{#if hasAnyPermission(user.permissions, ['user:create', 'situs:verify'])}
					<li>
						<a href="/dashboard" class="active:bg-primary/10 active:text-primary">
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-5 w-5"
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
							Dashboard
						</a>
					</li>
				{/if}

				<li class="menu-title">
					<span>Data Management</span>
				</li>

				<li>
					<a
						href="/dashboard/situs"
						class:hidden={!hasAnyPermission(user.permissions, ['situs:read_all', 'situs:read_own'])}
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-5 w-5"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
							/>
						</svg>
						Daftar Situs
					</a>
				</li>

				<li>
					<a
						href="/dashboard/jenis-situs"
						class:hidden={!hasAnyPermission(user.permissions, ['jenis-situs:read'])}
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-5 w-5"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"
							/>
						</svg>
						Jenis Situs
					</a>
				</li>

				<li>
					<a
						href="/dashboard/petugas"
						class:hidden={!hasAnyPermission(user.permissions, ['user:read'])}
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-5 w-5"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"
							/>
						</svg>
						Petugas
					</a>
				</li>

				<li class="menu-title">
					<span>System</span>
				</li>

				<li>
					<a
						href="/dashboard/policy"
						class:hidden={!hasAnyPermission(user.permissions, ['policy:read'])}
					>
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-5 w-5"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"
							/>
						</svg>
						Hak Akses
					</a>
				</li>

				<li>
					<a href="/dashboard/profile">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="h-5 w-5"
							fill="none"
							viewBox="0 0 24 24"
							stroke="currentColor"
						>
							<path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"
							/>
						</svg>
						Profile
					</a>
				</li>
			</ul>

			<div class="mt-auto border-t border-base-200 p-4">
				<div class="alert alert-info">
					<svg
						xmlns="http://www.w3.org/2000/svg"
						fill="none"
						viewBox="0 0 24 24"
						class="h-6 w-6 shrink-0 stroke-current"
					>
						<path
							stroke-linecap="round"
							stroke-linejoin="round"
							stroke-width="2"
							d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
						></path>
					</svg>
					<span>Versi 0.1.0</span>
				</div>
			</div>
		</aside>
	</div>
</div>
