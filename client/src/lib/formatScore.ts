export const formatScore = (score: number) => score.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ',')
