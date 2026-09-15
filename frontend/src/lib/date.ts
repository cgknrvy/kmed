/**
 * Calculates age from the give date of birth string. Date and Month
 * of year are taken into consideration.
 *
 * @param dob a valid date string representing the date of birth.
 */
export function calculateAge(dob: string): number {
  const dobDate = new Date(dob);
  const now = new Date();

  let years = now.getFullYear() - dobDate.getFullYear();

  // Check if the birth date and month has passed
  const hasNotHadAnniversary =
    now.getMonth() < dobDate.getMonth() ||
    (now.getMonth() === dobDate.getMonth() &&
      now.getDate() < dobDate.getDate());

  if (hasNotHadAnniversary) {
    years--;
  }

  return years;
}

export function parseDate(dateString: string): string {
  const date = new Date(dateString);
  return date.toLocaleString("en-KE");
}
