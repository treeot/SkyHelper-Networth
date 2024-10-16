package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Constants
var (
	applicationWorth = map[string]float64{
		"enrichment":            0.5,
		"farmingForDummies":     0.5,
		"gemstonePowerScroll":   0.5,
		"woodSingularity":       0.5,
		"artOfWar":              0.6,
		"fumingPotatoBook":      0.6,
		"gemstoneSlots":         0.6,
		"runes":                 0.6,
		"tunedTransmission":     0.7,
		"pocketSackInASack":     0.7,
		"essence":               0.75,
		"goldenBounty":          0.75,
		"silex":                 0.75,
		"artOfPeace":            0.8,
		"divanPowderCoating":    0.8,
		"jalapenoBook":          0.8,
		"manaDisintegrator":     0.8,
		"recomb":                0.8,
		"thunderInABottle":      0.8,
		"enchants":              0.85,
		"shensAuctionPrice":     0.85,
		"dye":                   0.9,
		"gemstoneChambers":      0.9,
		"attributes":            1,
		"drillPart":             1,
		"etherwarp":             1,
		"masterStar":            1,
		"gemstone":              1,
		"hotPotatoBook":         1,
		"necronBladeScroll":     1,
		"polarvoid":             1,
		"prestigeItem":          1,
		"reforge":               1,
		"winningBid":            1,
		"petCandy":              0.65,
		"soulboundPetSkins":     0.8,
		"petItem":               1,
	}
	enchantsWorth = map[string]float64{
		"counter_strike":      0.2,
		"big_brain":           0.35,
		"ultimate_inferno":    0.35,
		"overload":            0.35,
		"ultimate_soul_eater": 0.35,
		"ultimate_fatal_tempo":0.65,
	}
	blockedEnchants = map[string][]string{
		"bone_boomerang":        {"overload", "power", "ultimate_soul_eater"},
		"death_bow":             {"overload", "power", "ultimate_soul_eater"},
		"gardening_axe":         {"replenish"},
		"gardening_hoe":         {"replenish"},
		"advanced_gardening_axe":{"replenish"},
		"advanced_gardening_hoe":{"replenish"},
	}
	ignoredEnchants = map[string]int{
		"scavenger": 5,
	}
	stackingEnchants = []string{"expertise", "compact", "cultivating", "champion", "hecatomb", "toxophilite"}
	ignoreSilex = []string{"promising_spade"}
	masterStars = []string{"first_master_star", "second_master_star", "third_master_star", "fourth_master_star", "fifth_master_star"}
	validRunes = []string{
		"MUSIC_1", "MUSIC_2", "MUSIC_3", "MEOW_MUSIC_3", "ENCHANT_1", "ENCHANT_2", "ENCHANT_3", "GRAND_SEARING_3", "SPELLBOUND_3", "GRAND_FREEZING_3",
		"PRIMAL_FEAR_3", "GOLDEN_CARPET_3", "ICE_SKATES_3", "BARK_TUNES_3", "SMITTEN_3", "RAINY_DAY_3",
	}
	allowedRecombTypes = []string{"ACCESSORY", "NECKLACE", "GLOVES", "BRACELET", "BELT", "CLOAK"}
	allowedRecombIds = []string{
		"divan_helmet", "divan_chestplate", "divan_leggings", "divan_boots", "fermento_helmet", "fermento_chestplate", "fermento_leggings", "fermento_boots",
		"shadow_assassin_cloak", "starred_shadow_assassin_cloak",
	}
	attributesBaseCosts = map[string]string{
		"glowstone_gauntlet": "glowstone_gauntlet",
		"vanquished_glowstone_gauntlet": "glowstone_gauntlet",
		"blaze_belt": "blaze_belt",
		"vanquished_blaze_belt": "blaze_belt",
		"magma_necklace": "magma_necklace",
		"vanquished_magma_necklace": "magma_necklace",
		"magma_rod": "magma_rod",
		"inferno_rod": "magma_rod",
		"hellfire_rod": "magma_rod",
	}
	enrichments = []string{
		"TALISMAN_ENRICHMENT_CRITICAL_CHANCE", "TALISMAN_ENRICHMENT_CRITICAL_DAMAGE", "TALISMAN_ENRICHMENT_DEFENSE", "TALISMAN_ENRICHMENT_HEALTH",
		"TALISMAN_ENRICHMENT_INTELLIGENCE", "TALISMAN_ENRICHMENT_MAGIC_FIND", "TALISMAN_ENRICHMENT_WALK_SPEED", "TALISMAN_ENRICHMENT_STRENGTH",
		"TALISMAN_ENRICHMENT_ATTACK_SPEED", "TALISMAN_ENRICHMENT_FEROCITY", "TALISMAN_ENRICHMENT_SEA_CREATURE_CHANCE",
	}
	pickonimbusDurability = 5000
	specialEnchantmentMatches = map[string]string{
		"aiming": "Dragon Tracer", "counter_strike": "Counter-Strike", "pristine": "Prismatic", "turbo_cacti": "Turbo-Cacti", "turbo_cane": "Turbo-Cane",
		"turbo_carrot": "Turbo-Carrot", "turbo_cocoa": "Turbo-Cocoa", "turbo_melon": "Turbo-Melon", "turbo_mushrooms": "Turbo-Mushrooms", "turbo_potato": "Turbo-Potato",
		"turbo_pumpkin": "Turbo-Pumpkin", "turbo_warts": "Turbo-Warts", "turbo_wheat": "Turbo-Wheat", "ultimate_reiterate": "Ultimate Duplex", "ultimate_bobbin_time": "Ultimate Bobbin' Time",
	}
	prestiges = map[string][]string{
		"HOT_CRIMSON_CHESTPLATE": {"CRIMSON_CHESTPLATE"}, "HOT_CRIMSON_HELMET": {"CRIMSON_HELMET"}, "HOT_CRIMSON_LEGGINGS": {"CRIMSON_LEGGINGS"}, "HOT_CRIMSON_BOOTS": {"CRIMSON_BOOTS"},
		"BURNING_CRIMSON_CHESTPLATE": {"HOT_CRIMSON_CHESTPLATE", "CRIMSON_CHESTPLATE"}, "BURNING_CRIMSON_HELMET": {"HOT_CRIMSON_HELMET", "CRIMSON_HELMET"},
		"BURNING_CRIMSON_LEGGINGS": {"HOT_CRIMSON_LEGGINGS", "CRIMSON_LEGGINGS"}, "BURNING_CRIMSON_BOOTS": {"HOT_CRIMSON_BOOTS", "CRIMSON_BOOTS"},
		"FIERY_CRIMSON_CHESTPLATE": {"BURNING_CRIMSON_CHESTPLATE", "HOT_CRIMSON_CHESTPLATE", "CRIMSON_CHESTPLATE"}, "FIERY_CRIMSON_HELMET": {"BURNING_CRIMSON_HELMET", "HOT_CRIMSON_HELMET", "CRIMSON_HELMET"},
		"FIERY_CRIMSON_LEGGINGS": {"BURNING_CRIMSON_LEGGINGS", "HOT_CRIMSON_LEGGINGS", "CRIMSON_LEGGINGS"}, "FIERY_CRIMSON_BOOTS": {"BURNING_CRIMSON_BOOTS", "HOT_CRIMSON_BOOTS", "CRIMSON_BOOTS"},
		"INFERNAL_CRIMSON_CHESTPLATE": {"FIERY_CRIMSON_CHESTPLATE", "BURNING_CRIMSON_CHESTPLATE", "HOT_CRIMSON_CHESTPLATE", "CRIMSON_CHESTPLATE"},
		"INFERNAL_CRIMSON_HELMET": {"FIERY_CRIMSON_HELMET", "BURNING_CRIMSON_HELMET", "HOT_CRIMSON_HELMET", "CRIMSON_HELMET"},
		"INFERNAL_CRIMSON_LEGGINGS": {"FIERY_CRIMSON_LEGGINGS", "BURNING_CRIMSON_LEGGINGS", "HOT_CRIMSON_LEGGINGS", "CRIMSON_LEGGINGS"},
		"INFERNAL_CRIMSON_BOOTS": {"FIERY_CRIMSON_BOOTS", "BURNING_CRIMSON_BOOTS", "HOT_CRIMSON_BOOTS", "CRIMSON_BOOTS"},
		"HOT_TERROR_CHESTPLATE": {"TERROR_CHESTPLATE"}, "HOT_TERROR_HELMET": {"TERROR_HELMET"}, "HOT_TERROR_LEGGINGS": {"TERROR_LEGGINGS"}, "HOT_TERROR_BOOTS": {"TERROR_BOOTS"},
		"BURNING_TERROR_CHESTPLATE": {"HOT_TERROR_CHESTPLATE", "TERROR_CHESTPLATE"}, "BURNING_TERROR_HELMET": {"HOT_TERROR_HELMET", "TERROR_HELMET"},
		"BURNING_TERROR_LEGGINGS": {"HOT_TERROR_LEGGINGS", "TERROR_LEGGINGS"}, "BURNING_TERROR_BOOTS": {"HOT_TERROR_BOOTS", "TERROR_BOOTS"},
		"FIERY_TERROR_CHESTPLATE": {"BURNING_TERROR_CHESTPLATE", "HOT_TERROR_CHESTPLATE", "TERROR_CHESTPLATE"}, "FIERY_TERROR_HELMET": {"BURNING_TERROR_HELMET", "HOT_TERROR_HELMET", "TERROR_HELMET"},
		"FIERY_TERROR_LEGGINGS": {"BURNING_TERROR_LEGGINGS", "HOT_TERROR_LEGGINGS", "TERROR_LEGGINGS"}, "FIERY_TERROR_BOOTS": {"BURNING_TERROR_BOOTS", "HOT_TERROR_BOOTS", "TERROR_BOOTS"},
		"INFERNAL_TERROR_CHESTPLATE": {"FIERY_TERROR_CHESTPLATE", "BURNING_TERROR_CHESTPLATE", "HOT_TERROR_CHESTPLATE", "TERROR_CHESTPLATE"},
		"INFERNAL_TERROR_HELMET": {"FIERY_TERROR_HELMET", "BURNING_TERROR_HELMET", "HOT_TERROR_HELMET", "TERROR_HELMET"},
		"INFERNAL_TERROR_LEGGINGS": {"FIERY_TERROR_LEGGINGS", "BURNING_TERROR_LEGGINGS", "HOT_TERROR_LEGGINGS", "TERROR_LEGGINGS"},
		"INFERNAL_TERROR_BOOTS": {"FIERY_TERROR_BOOTS", "BURNING_TERROR_BOOTS", "HOT_TERROR_BOOTS", "TERROR_BOOTS"},
		"HOT_FERVOR_CHESTPLATE": {"FERVOR_CHESTPLATE"}, "HOT_FERVOR_HELMET": {"FERVOR_HELMET"}, "HOT_FERVOR_LEGGINGS": {"FERVOR_LEGGINGS"}, "HOT_FERVOR_BOOTS": {"FERVOR_BOOTS"},
		"BURNING_FERVOR_CHESTPLATE": {"HOT_FERVOR_CHESTPLATE", "FERVOR_CHESTPLATE"}, "BURNING_FERVOR_HELMET": {"HOT_FERVOR_HELMET", "FERVOR_HELMET"},
		"BURNING_FERVOR_LEGGINGS": {"HOT_FERVOR_LEGGINGS", "FERVOR_LEGGINGS"}, "BURNING_FERVOR_BOOTS": {"HOT_FERVOR_BOOTS", "FERVOR_BOOTS"},
		"FIERY_FERVOR_CHESTPLATE": {"BURNING_FERVOR_CHESTPLATE", "HOT_FERVOR_CHESTPLATE", "FERVOR_CHESTPLATE"}, "FIERY_FERVOR_HELMET": {"BURNING_FERVOR_HELMET", "HOT_FERVOR_HELMET", "FERVOR_HELMET"},
		"FIERY_FERVOR_LEGGINGS": {"BURNING_FERVOR_LEGGINGS", "HOT_FERVOR_LEGGINGS", "FERVOR_LEGGINGS"}, "FIERY_FERVOR_BOOTS": {"BURNING_FERVOR_BOOTS", "HOT_FERVOR_BOOTS", "FERVOR_BOOTS"},
		"INFERNAL_FERVOR_CHESTPLATE": {"FIERY_FERVOR_CHESTPLATE", "BURNING_FERVOR_CHESTPLATE", "HOT_FERVOR_CHESTPLATE", "FERVOR_CHESTPLATE"},
		"INFERNAL_FERVOR_HELMET": {"FIERY_FERVOR_HELMET", "BURNING_FERVOR_HELMET", "HOT_FERVOR_HELMET", "FERVOR_HELMET"},
		"INFERNAL_FERVOR_LEGGINGS": {"FIERY_FERVOR_LEGGINGS", "BURNING_FERVOR_LEGGINGS", "HOT_FERVOR_LEGGINGS", "FERVOR_LEGGINGS"},
		"INFERNAL_FERVOR_BOOTS": {"FIERY_FERVOR_BOOTS", "BURNING_FERVOR_BOOTS", "HOT_FERVOR_BOOTS", "FERVOR_BOOTS"},
		"HOT_HOLLOW_CHESTPLATE": {"HOLLOW_CHESTPLATE"}, "HOT_HOLLOW_HELMET": {"HOLLOW_HELMET"}, "HOT_HOLLOW_LEGGINGS": {"HOLLOW_LEGGINGS"}, "HOT_HOLLOW_BOOTS": {"HOLLOW_BOOTS"},
		"BURNING_HOLLOW_CHESTPLATE": {"HOT_HOLLOW_CHESTPLATE", "HOLLOW_CHESTPLATE"}, "BURNING_HOLLOW_HELMET": {"HOT_HOLLOW_HELMET", "HOLLOW_HELMET"},
		"BURNING_HOLLOW_LEGGINGS": {"HOT_HOLLOW_LEGGINGS", "HOLLOW_LEGGINGS"}, "BURNING_HOLLOW_BOOTS": {"HOT_HOLLOW_BOOTS", "HOLLOW_BOOTS"},
		"FIERY_HOLLOW_CHESTPLATE": {"BURNING_HOLLOW_CHESTPLATE", "HOT_HOLLOW_CHESTPLATE", "HOLLOW_CHESTPLATE"}, "FIERY_HOLLOW_HELMET": {"BURNING_HOLLOW_HELMET", "HOT_HOLLOW_HELMET", "HOLLOW_HELMET"},
		"FIERY_HOLLOW_LEGGINGS": {"BURNING_HOLLOW_LEGGINGS", "HOT_HOLLOW_LEGGINGS", "HOLLOW_LEGGINGS"}, "FIERY_HOLLOW_BOOTS": {"BURNING_HOLLOW_BOOTS", "HOT_HOLLOW_BOOTS", "HOLLOW_BOOTS"},
		"INFERNAL_HOLLOW_CHESTPLATE": {"FIERY_HOLLOW_CHESTPLATE", "BURNING_HOLLOW_CHESTPLATE", "HOT_HOLLOW_CHESTPLATE", "HOLLOW_CHESTPLATE"},
		"INFERNAL_HOLLOW_HELMET": {"FIERY_HOLLOW_HELMET", "BURNING_HOLLOW_HELMET", "HOT_HOLLOW_HELMET", "HOLLOW_HELMET"},
		"INFERNAL_HOLLOW_LEGGINGS": {"FIERY_HOLLOW_LEGGINGS", "BURNING_HOLLOW_LEGGINGS", "HOT_HOLLOW_LEGGINGS", "HOLLOW_LEGGINGS"},
		"INFERNAL_HOLLOW_BOOTS": {"FIERY_HOLLOW_BOOTS", "BURNING_HOLLOW_BOOTS", "HOT_HOLLOW_BOOTS", "HOLLOW_BOOTS"},
		"HOT_AURORA_CHESTPLATE": {"AURORA_CHESTPLATE"}, "HOT_AURORA_HELMET": {"AURORA_HELMET"}, "HOT_AURORA_LEGGINGS": {"AURORA_LEGGINGS"}, "HOT_AURORA_BOOTS": {"AURORA_BOOTS"},
		"BURNING_AURORA_CHESTPLATE": {"HOT_AURORA_CHESTPLATE", "AURORA_CHESTPLATE"}, "BURNING_AURORA_HELMET": {"HOT_AURORA_HELMET", "AURORA_HELMET"},
		"BURNING_AURORA_LEGGINGS": {"HOT_AURORA_LEGGINGS", "AURORA_LEGGINGS"}, "BURNING_AURORA_BOOTS": {"HOT_AURORA_BOOTS", "AURORA_BOOTS"},
		"FIERY_AURORA_CHESTPLATE": {"BURNING_AURORA_CHESTPLATE", "HOT_AURORA_CHESTPLATE", "AURORA_CHESTPLATE"}, "FIERY_AURORA_HELMET": {"BURNING_AURORA_HELMET", "HOT_AURORA_HELMET", "AURORA_HELMET"},
		"FIERY_AURORA_LEGGINGS": {"BURNING_AURORA_LEGGINGS", "HOT_AURORA_LEGGINGS", "AURORA_LEGGINGS"}, "FIERY_AURORA_BOOTS": {"BURNING_AURORA_BOOTS", "HOT_AURORA_BOOTS", "AURORA_BOOTS"},
		"INFERNAL_AURORA_CHESTPLATE": {"FIERY_AURORA_CHESTPLATE", "BURNING_AURORA_CHESTPLATE", "HOT_AURORA_CHESTPLATE", "AURORA_CHESTPLATE"},
		"INFERNAL_AURORA_HELMET": {"FIERY_AURORA_HELMET", "BURNING_AURORA_HELMET", "HOT_AURORA_HELMET", "AURORA_HELMET"},
		"INFERNAL_AURORA_LEGGINGS": {"FIERY_AURORA_LEGGINGS", "BURNING_AURORA_LEGGINGS", "HOT_AURORA_LEGGINGS", "AURORA_LEGGINGS"},
		"INFERNAL_AURORA_BOOTS": {"FIERY_AURORA_BOOTS", "BURNING_AURORA_BOOTS", "HOT_AURORA_BOOTS", "AURORA_BOOTS"},
	}
	reforges = map[string]string{
		"stiff": "hardened_wood", "salty": "salt_cube", "aote_stone": "aote_stone", "blazing": "blazen_sphere", "waxed": "blaze_wax", "rooted": "burrowing_spores",
		"candied": "candy_corn", "perfect": "diamond_atom", "fleet": "diamonite", "fabled": "dragon_claw", "spiked": "dragon_scale", "royal": "dwarven_treasure",
		"hyper": "endstone_geode", "coldfusion": "entropy_suppressor", "blooming": "flowering_bouquet", "fanged": "full_jaw_fanging_kit", "jaded": "jaderald",
		"jerry": "jerry_stone", "magnetic": "lapis_crystal", "earthy": "large_walnut", "fortified": "meteor_shard", "gilded": "midas_jewel", "cubic": "molten_cube",
		"necrotic": "necromancer_brooch", "fruitful": "onyx", "precise": "optical_lens", "mossy": "overgrown_grass", "pitchin": "pitchin_koi", "undead": "premium_flesh",
		"blood_soaked": "presumed_gallon_of_red_paint", "mithraic": "pure_mithril", "reinforced": "rare_diamond", "ridiculous": "red_nose", "loving": "red_scarf",
		"auspicious": "rock_gemstone", "treacherous": "rusty_anchor", "headstrong": "salmon_opal", "strengthened": "searing_stone", "glistening": "shiny_prism",
		"bustling": "skymart_brochure", "spiritual": "spirit_decoy", "suspicious": "suspicious_vial", "snowy": "terry_snowglobe", "dimensional": "titanium_tesseract",
		"ambered": "amber_material", "beady": "beady_eyes", "blessed": "blessed_fruit", "bulky": "bulky_stone", "buzzing": "clipped_wings", "submerged": "deep_sea_orb",
		"renowned": "dragon_horn", "festive": "frozen_bauble", "giant": "giant_tooth", "lustrous": "gleaming_crystal", "bountiful": "golden_ball", "chomp": "kuudra_mandible",
		"lucky": "lucky_dice", "stellar": "petrified_starfall", "scraped": "pocket_iceberg", "ancient": "precursor_gear", "refined": "refined_amber", "empowered": "sadan_brooch",
		"withered": "wither_blood", "glacial": "frigid_husk", "heated": "hot_stuff", "dirty": "dirt_bottle", "moil": "moil_log", "toil": "toil_log", "greater_spook": "boo_stone",
	}
)

// Helper functions
func titleCase(str string) string {
	words := strings.Split(strings.ToLower(str), "_")
	for i, word := range words {
		words[i] = strings.Title(word)
	}
	return strings.Join(words, " ")
}

func decodeData(data string) ([]interface{}, error) {
	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}
	var result []interface{}
	err = json.Unmarshal(decoded, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Calculators
func calculateEssence(item map[string]interface{}, prices map[string]float64) map[string]interface{} {
	itemID := strings.ToLower(item["id"].(string))
	itemPrice := prices[itemID]
	if itemPrice > 0 {
		return map[string]interface{}{
			"name":        fmt.Sprintf("%s Essence", titleCase(strings.Split(itemID, "_")[1])),
			"id":          itemID,
			"price":       itemPrice * item["amount"].(float64),
			"calculation": []interface{}{},
			"count":       item["amount"],
			"soulbound":   false,
		}
	}
	return nil
}

func calculateItem(item map[string]interface{}, prices map[string]float64, returnItemData bool) map[string]interface{} {
	// TODO: Implement Backpack Calculations

	if item["tag"].(map[string]interface{})["ExtraAttributes"].(map[string]interface{})["id"] == "PET" && item["tag"].(map[string]interface{})["ExtraAttributes"].(map[string]interface{})["petInfo"] != nil {
		petInfo := item["tag"].(map[string]interface{})["ExtraAttributes"].(map[string]interface{})["petInfo"]
		if petInfoStr, ok := petInfo.(string); ok {
			json.Unmarshal([]byte(petInfoStr), &petInfo)
		}
		level := getPetLevel(petInfo.(map[string]interface{}))
		petInfo.(map[string]interface{})["level"] = level["level"]
		petInfo.(map[string]interface{})["xpMax"] = level["xpMax"]
		return calculatePet(petInfo.(map[string]interface{}), prices)
	}

	if item["tag"].(map[string]interface{})["ExtraAttributes"].(map[string]interface{})["id"] != nil {
		if item["tag"].(map[string]interface{})["display"] == nil {
			return nil
		}
		itemName := strings.ReplaceAll(item["tag"].(map[string]interface{})["display"].(map[string]interface{})["Name"].(string), "§[0-9a-fk-or]", "")
		itemID := strings.ToLower(item["tag"].(map[string]interface{})["ExtraAttributes"].(map[string]interface{})["id"].(string))
		ExtraAttributes := item["tag"].(map[string]interface{})["ExtraAttributes"].(map[string]interface{})
		skyblockItem := getHypixelItemInformationFromId(strings.ToUpper(itemID))

		if ExtraAttributes["skin"] != nil {
			if prices[fmt.Sprintf("%s_skinned_%s", itemID, strings.ToLower(ExtraAttributes["skin"].(string)))] > 0 {
				itemID = fmt.Sprintf("%s_skinned_%s", itemID, strings.ToLower(ExtraAttributes["skin"].(string)))
			}
		}
		if itemID == "party_hat_sloth" && ExtraAttributes["party_hat_emoji"] != nil {
			if prices[fmt.Sprintf("%s_%s", itemID, strings.ToLower(ExtraAttributes["party_hat_emoji"].(string)))] > 0 {
				itemID = fmt.Sprintf("%s_%s", itemID, strings.ToLower(ExtraAttributes["party_hat_emoji"].(string)))
			}
		}

		if itemName == "Beastmaster Crest" || itemName == "Griffin Upgrade Stone" || itemName == "Wisp Upgrade Stone" {
			itemName = fmt.Sprintf("%s (%s)", itemName, titleCase(strings.ReplaceAll(skyblockItem["tier"].(string), "_", " ")))
		} else if strings.HasSuffix(itemName, " Exp Boost") {
			itemName = fmt.Sprintf("%s (%s)", itemName, titleCase(strings.Split(skyblockItem["id"].(string), "_")[len(strings.Split(skyblockItem["id"].(string), "_"))-1]))
		}

		if (ExtraAttributes["id"] == "RUNE" || ExtraAttributes["id"] == "UNIQUE_RUNE") && ExtraAttributes["runes"] != nil && len(ExtraAttributes["runes"].(map[string]interface{})) > 0 {
			runeType := ""
			runeTier := ""
			for k, v := range ExtraAttributes["runes"].(map[string]interface{}) {
				runeType = k
				runeTier = v.(string)
				break
			}
			itemID = fmt.Sprintf("rune_%s_%s", runeType, runeTier)
		}
		if ExtraAttributes["id"] == "NEW_YEAR_CAKE" {
			itemID = fmt.Sprintf("new_year_cake_%s", ExtraAttributes["new_years_cake"].(string))
		}
		if ExtraAttributes["id"] == "PARTY_HAT_CRAB" || ExtraAttributes["id"] == "PARTY_HAT_CRAB_ANIMATED" || ExtraAttributes["id"] == "BALLOON_HAT_2024" {
			if ExtraAttributes["party_hat_color"] != nil {
				itemID = fmt.Sprintf("%s_%s", strings.ToLower(ExtraAttributes["id"].(string)), ExtraAttributes["party_hat_color"].(string))
			}
		}
		if ExtraAttributes["id"] == "DCTR_SPACE_HELM" && ExtraAttributes["edition"] != nil {
			itemID = "dctr_space_helm_editioned"
		}
		if ExtraAttributes["id"] == "CREATIVE_MIND" && ExtraAttributes["edition"] == nil {
			itemID = "creative_mind_uneditioned"
		}
		if ExtraAttributes["is_shiny"] != nil && prices[fmt.Sprintf("%s_shiny", itemID)] > 0 {
			itemID = fmt.Sprintf("%s_shiny", itemID)
		}
		if strings.HasPrefix(ExtraAttributes["id"].(string), "STARRED_") && prices[itemID] == 0 && prices[strings.Replace(itemID, "starred_", "", 1)] > 0 {
			itemID = strings.Replace(itemID, "starred_", "", 1)
		}

		itemData := prices[itemID]
		price := itemData * item["Count"].(float64)
		base := itemData * item["Count"].(float64)
		if ExtraAttributes["skin"] != nil {
			newPrice := prices[strings.ToLower(item["tag"].(map[string]interface{})["ExtraAttributes"].(map[string]interface{})["id"].(string))]
			if newPrice > price {
				price = newPrice * item["Count"].(float64)
				base = newPrice * item["Count"].(float64)
			}
		}
		if price == 0 && ExtraAttributes["price"] != nil {
			price = float64(ExtraAttributes["price"].(int)) * 0.85
			base = float64(ExtraAttributes["price"].(int)) * 0.85
		}
		calculation := []interface{}{}

		if ExtraAttributes["id"] == "PICKONIMBUS" && ExtraAttributes["pickonimbus_durability"] != nil {
			reduction := float64(ExtraAttributes["pickonimbus_durability"].(int)) / float64(pickonimbusDurability)
			price += price * (reduction - 1)
			base += price * (reduction - 1)
		}

		if itemID != "attribute_shard" && ExtraAttributes["attributes"] != nil {
			sortedAttributes := make([]string, 0, len(ExtraAttributes["attributes"].(map[string]interface{})))
			for k := range ExtraAttributes["attributes"].(map[string]interface{}) {
				sortedAttributes = append(sortedAttributes, k)
			}
			sort.Strings(sortedAttributes)
			formattedID := strings.ReplaceAll(itemID, "(hot_|fiery_|burning_|infernal_)", "")
			godRollID := fmt.Sprintf("%s%s", formattedID, strings.Join(sortedAttributes, "_roll_"))
			godRollPrice := prices[godRollID]
			if godRollPrice > price {
				price = godRollPrice
				base = godRollPrice
				calculation = append(calculation, map[string]interface{}{
					"id":    godRollID[len(formattedID)+1:],
					"type":  "god_roll",
					"price": godRollPrice,
					"count": 1,
				})
			}
		}

		if itemData == 0 {
			prestige := prestiges[strings.ToUpper(itemID)]
			if prestige != nil {
				for _, prestigeItem := range prestige {
					foundItem := getHypixelItemInformationFromId(prestigeItem)
					if price == 0 {
						price = 0
					}
					if foundItem["upgrade_costs"] != nil {
						price += starCosts(prices, calculation, foundItem["upgrade_costs"].([]interface{}), prestigeItem)
					}
					if foundItem["prestige"] != nil && foundItem["prestige"].(map[string]interface{})["costs"] != nil {
						price += starCosts(prices, calculation, foundItem["prestige"].(map[string]interface{})["costs"].([]interface{}), prestigeItem)
					}
				}
			}
		}

		if ExtraAttributes["price"] != nil && ExtraAttributes["auction"] != nil && ExtraAttributes["bid"] != nil {
			pricePaid := float64(ExtraAttributes["price"].(int)) * applicationWorth["shensAuctionPrice"]
			if pricePaid > price {
				price = pricePaid
				calculation = append(calculation, map[string]interface{}{
					"id":    itemID,
					"type":  "shens_auction",
					"price": pricePaid,
					"count": 1,
				})
			}
		}

		if itemID == "enchanted_book" && ExtraAttributes["enchantments"] != nil {
			if len(ExtraAttributes["enchantments"].(map[string]interface{})) == 1 {
				for name, value := range ExtraAttributes["enchantments"].(map[string]interface{}) {
					calculation = append(calculation, map[string]interface{}{
						"id":    fmt.Sprintf("%s_%d", strings.ToUpper(name), value.(int)),
						"type":  "enchant",
						"price": prices[fmt.Sprintf("enchantment_%s_%d", strings.ToLower(name), value.(int))],
						"count": 1,
					})
					price = prices[fmt.Sprintf("enchantment_%s_%d", strings.ToLower(name), value.(int))]
					itemName = specialEnchantmentMatches[name]
				}
			} else {
				enchantmentPrice := 0.0
				for name, value := range ExtraAttributes["enchantments"].(map[string]interface{}) {
					calculation = append(calculation, map[string]interface{}{
						"id":    fmt.Sprintf("%s_%d", strings.ToUpper(name), value.(int)),
						"type":  "enchant",
						"price": prices[fmt.Sprintf("enchantment_%s_%d", strings.ToLower(name), value.(int))] * applicationWorth["enchants"],
						"count": 1,
					})
					enchantmentPrice += prices[fmt.Sprintf("enchantment_%s_%d", strings.ToLower(name), value.(int))] * applicationWorth["enchants"]
				}
				price = enchantmentPrice
			}
		} else if ExtraAttributes["enchantments"] != nil {
			for name, value := range ExtraAttributes["enchantments"].(map[string]interface{}) {
				name = strings.ToLower(name)
				if contains(blockedEnchants[itemID], name) {
					continue
				}
				if ignoredEnchants[name] == value {
					continue
				}
				if contains(stackingEnchants, name) {
					value = 1
				}
				if name == "efficiency" && value.(int) > 5 && !contains(ignoreSilex, itemID) {
					efficiencyLevel := value.(int) - 5
					if itemID == "stonk_pickaxe" {
						efficiencyLevel--
					}
					if efficiencyLevel > 0 {
						calculation = append(calculation, map[string]interface{}{
							"id":    "SIL_EX",
							"type":  "silex",
							"price": prices["sil_ex"] * float64(efficiencyLevel) * applicationWorth["silex"],
							"count": efficiencyLevel,
						})
						price += prices["sil_ex"] * float64(efficiencyLevel) * applicationWorth["silex"]
					}
				}
				if name == "scavenger" && value.(int) >= 6 {
					calculation = append(calculation, map[string]interface{}{
						"id":    "GOLDEN_BOUNTY",
						"type":  "golden_bounty",
						"price": prices["GOLDEN_BOUNTY"] * applicationWorth["goldenBounty"],
						"count": 1,
					})
					price += prices["GOLDEN_BOUNTY"] * applicationWorth["goldenBounty"]
				}
				calculation = append(calculation, map[string]interface{}{
					"id":    fmt.Sprintf("%s_%d", strings.ToUpper(name), value.(int)),
					"type":  "enchant",
					"price": prices[fmt.Sprintf("enchantment_%s_%d", name, value.(int))] * enchantsWorth[name],
					"count": 1,
				})
				price += prices[fmt.Sprintf("enchantment_%s_%d", name, value.(int))] * enchantsWorth[name]
			}
		}

		if ExtraAttributes["attributes"] != nil {
			for attribute, tier := range ExtraAttributes["attributes"].(map[string]interface{}) {
				if tier.(int) == 1 {
					continue
				}
				shards := (1 << (tier.(int) - 1)) - 1
				baseAttributePrice := prices[fmt.Sprintf("attribute_shard_%s", attribute)]
				if attributesBaseCosts[itemID] != "" && prices[attributesBaseCosts[itemID]] < baseAttributePrice {
					baseAttributePrice = prices[attributesBaseCosts[itemID]]
				} else if strings.HasPrefix(itemID, "aurora") && prices[fmt.Sprintf("kuudra_helmet_%s", attribute)] < baseAttributePrice {
					baseAttributePrice = prices[fmt.Sprintf("kuudra_helmet_%s", attribute)]
				} else if strings.HasPrefix(itemID, "aurora") {
					kuudraPrices := []float64{prices[fmt.Sprintf("kuudra_chestplate_%s", attribute)], prices[fmt.Sprintf("kuudra_leggings_%s", attribute)], prices[fmt.Sprintf("kuudra_boots_%s", attribute)]}
					kuudraPrice := 0.0
					for _, v := range kuudraPrices {
						kuudraPrice += v
					}
					kuudraPrice /= float64(len(kuudraPrices))
					if kuudraPrice > 0 && (baseAttributePrice == 0 || kuudraPrice < baseAttributePrice) {
						baseAttributePrice = kuudraPrice
					}
				}
				if baseAttributePrice == 0 {
					continue
				}
				attributePrice := baseAttributePrice * float64(shards) * applicationWorth["attributes"]
				price += attributePrice
				calculation = append(calculation, map[string]interface{}{
					"id":      fmt.Sprintf("%s_%d", strings.ToUpper(attribute), tier.(int)),
					"type":    "attribute",
					"price":   attributePrice,
					"count":   1,
					"shards":  shards,
				})
			}
		}

		if ExtraAttributes["sack_pss"] != nil {
			calculation = append(calculation, map[string]interface{}{
				"id":    "POCKET_SACK_IN_A_SACK",
				"type":  "pocket_sack_in_a_sack",
				"price": prices["pocket_sack_in_a_sack"] * float64(ExtraAttributes["sack_pss"].(int)) * applicationWorth["pocketSackInASack"],
				"count": ExtraAttributes["sack_pss"],
			})
			price += prices["pocket_sack_in_a_sack"] * float64(ExtraAttributes["sack_pss"].(int)) * applicationWorth["pocketSackInASack"]
		}

		if ExtraAttributes["wood_singularity_count"] != nil {
			calculation = append(calculation, map[string]interface{}{
				"id":    "WOOD_SINGULARITY",
				"type":  "wood_singularity",
				"price": prices["wood_singularity"] * float64(ExtraAttributes["wood_singularity_count"].(int)) * applicationWorth["woodSingularity"],
				"count": ExtraAttributes["wood_singularity_count"],
			})
			price += prices["wood_singularity"] * float64(ExtraAttributes["wood_singularity_count"].(int)) * applicationWorth["woodSingularity"]
		}

		if ExtraAttributes["jalapeno_count"] != nil {
			calculation = append(calculation, map[string]interface{}{
				"id":    "JALAPENO_BOOK",
				"type":  "jalapeno_book",
				"price": prices["jalapeno_book"] * float64(ExtraAttributes["jalapeno_count"].(int)) * applicationWorth["jalapenoBook"],
				"count": ExtraAttributes["jalapeno_count"],
			})
			price += prices["jalapeno_book"] * float64(ExtraAttributes["jalapeno_count"].(int)) * applicationWorth["jalapenoBook"]
		}

		if ExtraAttributes["tuned_transmission"] != nil {
			calculation = append(calculation, map[string]interface{}{
				"id":    "TRANSMISSION_TUNER",
				"type":  "tuned_transmission",
				"price": prices["transmission_tuner"] * float64(ExtraAttributes["tuned_transmission"].(int)) * applicationWorth["tunedTransmission"],
				"count": ExtraAttributes["tuned_transmission"],
			})
			price += prices["transmission_tuner"] * float64(ExtraAttributes["tuned_transmission"].(int)) * applicationWorth["tunedTransmission"]
		}

		if ExtraAttributes["mana_disintegrator_count"] != nil {
			calculation = append(calculation, map[string]interface{}{
				"id":    "MANA_DISINTEGRATOR",
				"type":  "mana_disintegrator",
				"price": prices["mana_disintegrator"] * float64(ExtraAttributes["mana_disintegrator_count"].(int)) * applicationWorth["manaDisintegrator"],
				"count": ExtraAttributes["mana_disintegrator_count"],
			})
			price += prices["mana_disintegrator"] * float64(ExtraAttributes["mana_disintegrator_count"].(int)) * applicationWorth["manaDisintegrator"]
		}

		if ExtraAttributes["thunder_charge"] != nil && itemID == "pulse_ring
