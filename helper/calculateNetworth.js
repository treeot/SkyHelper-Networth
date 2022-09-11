const { calculatePet } = require('../calculators/petCalculator');
const { calculateSackItem } = require('../calculators/sacksCalculator');
const { calculateItem } = require('../calculators/itemCalculator');

const calculateNetworth = async (items, purseBalance, bankBalance, prices, onlyNetworth) => {
  const categories = {};

  for (const [category, categoryItems] of Object.entries(items)) {
    // Calculate networth for each category
    categories[category] = { total: 0, unsoulboundTotal: 0, items: [] };

    for (const item of categoryItems) {
      const result = category === 'pets' ? calculatePet(item, prices) : category === 'sacks' ? calculateSackItem(item, prices) : calculateItem(item, prices);

      categories[category].total += result?.price || 0;
      if (!result?.soulbound) categories[category].unsoulboundTotal += result?.price || 0;
      if (!onlyNetworth && result) categories[category].items.push(result);
    }

    // Sort items by price
    for (const category of Object.keys(categories)) {
      if (!onlyNetworth) categories[category].items.sort((a, b) => b.price - a.price);
    }

    if (onlyNetworth) delete categories[category].items;
  }

  // Calculate total networth
  const total = Object.values(categories).reduce((acc, category) => acc + category.total, 0) + (bankBalance || 0) + (purseBalance || 0);
  const unsoulboundTotal = Object.values(categories).reduce((acc, category) => acc + category.unsoulboundTotal, 0) + (bankBalance || 0) + (purseBalance || 0);

  return {
    noInventory: !items.inventory?.length,
    networth: total,
    unsoulboundNetworth: unsoulboundTotal,
    purse: purseBalance || 0,
    bank: bankBalance || 0,
    types: categories,
  };
};

module.exports = {
  calculateNetworth,
};
