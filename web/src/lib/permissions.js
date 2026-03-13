/**
 * Mengecek apakah user memiliki MINIMAL SATU dari permission yang dibutuhkan (Logika OR)
 * * @param {string[]} userPermissions - Array permission milik user (contoh: data.user.permissions)
 * @param {string | string[]} requiredPermissions - Permission yang mau dicek
 * @returns {boolean}
 */
export function hasAnyPermission(userPermissions, requiredPermissions) {
	// Return false jika data userPermissions kosong/tidak valid
	if (!userPermissions || !Array.isArray(userPermissions)) return false;

	// Normalisasi requiredPermissions menjadi array jika formatnya string tunggal
	const required = Array.isArray(requiredPermissions) ? requiredPermissions : [requiredPermissions];

	// Cek apakah ada minimal 1 permission yang cocok
	return required.some((permission) => userPermissions.includes(permission));
}

/**
 * Mengecek apakah user memiliki SEMUA permission yang dibutuhkan (Logika AND)
 * * @param {string[]} userPermissions - Array permission milik user
 * @param {string | string[]} requiredPermissions - Permission yang mau dicek
 * @returns {boolean}
 */
export function hasAllPermissions(userPermissions, requiredPermissions) {
	if (!userPermissions || !Array.isArray(userPermissions)) return false;

	const required = Array.isArray(requiredPermissions) ? requiredPermissions : [requiredPermissions];

	// Cek apakah SEMUA permission cocok
	return required.every((permission) => userPermissions.includes(permission));
}
