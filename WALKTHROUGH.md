# Walkthrough: Stored Races Wave Seeding & Season Qualifications

## Summary of Completed Work
We integrated Birkie and Kortelopet wave seeding brackets, qualifier race adjustments, and season-based wave qualification into the Go & React race tracker application.

The application automatically evaluates stored races (`PercentBack`, `RaceDate`, and `RaceName`) in MongoDB and calculates exact wave placements by ski season for:
1. **American Birkebeiner Skate 50K** (Waves 1 to 8)
2. **American Birkebeiner Classic 53K** (Waves 1 to 5)
3. **Kortelopet Skate 29K** (Waves 1 to 3)
4. **Kortelopet Classic 29K** (Waves 1 to 3)

---

## Qualifier Updates from Official Birkie Qualifier Table
All qualifier races and adjustment factors were updated directly from [birkie.com/ski/events/qualifier-races/](https://www.birkie.com/ski/events/qualifier-races/):
- **Mora Vasaloppet 34K Skate:** **30%** (previously placeholder 5%)
- **Mora Vasaloppet 51K Skate:** **13%**
- **Mora Vasaloppet 42K Classic:** **36%**
- **Giant's Ridge Pepsi Challenge 50K Skate:** **19%** (matches "Challenge Skate")
- **Giant's Ridge Pepsi Challenge 50K Classic:** **32%**
- **Giant's Ridge Pepsi Challenge 25K (Skate & Classic):** **30%**
- **Seeley Hills Classic 35K:** **30%**
- **Seeley Hills Classic 20K:** **57%**
- **Hayward Lion's Pre-Birkie 42K Skate:** **35%** (Classic: **26%**)
- **City of Lakes Loppet 31K Skate & 29K Classic:** **30%**
- **Noquemanon Ski Marathon 50K Skate:** **15%** (Classic: **14%**)
- **Great Bear Chase 50K Skate:** **29%** (Classic: **50%**)
- **Sisu Ski Fest 30K Skate:** **18%** (Classic: **20%**)
- **Craftsbury Marathon:** 30K Sk: **12%**, 45K Cl: **10%**, 30K Cl: **16%**
- Plus all other official regional marathons (Wolf Tracks, Boulder Mountain, North End, Vasa, etc.).
- Direct Birkie events (Birkie Skate, Birkie Classic) have **0% adjustment** for their direct event, with technique transfers applied only when crossing disciplines.

---

## Verification & Validation
- **Go Backend Build:** `go build ./...` passed with zero errors.
- **Frontend Production Build:** `npm run build` passed with zero warnings.
- **Static Assets:** Deployed to `./static/` for immediate use.
