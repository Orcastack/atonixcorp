package core

import (
	"fmt"
)

type PricingEngine struct {
	Plans         map[string]Plan
	RegionPricing map[string]RegionPrice
	EnginePricing map[string]float64
}

type RegionPrice struct {
	Multiplier float64 `yaml:"multiplier"`
}

func NewPricingEngine(plans map[string]Plan, rp map[string]RegionPrice, ep map[string]float64) *PricingEngine {
	return &PricingEngine{
		Plans:         plans,
		RegionPricing: rp,
		EnginePricing: ep,
	}
}

func (p *PricingEngine) Calculate(inst *DBInstance) (float64, error) {
	plan, ok := p.Plans[inst.Plan]
	if !ok {
		return 0, fmt.Errorf("invalid plan: %s", inst.Plan)
	}

	// Base plan price
	total := float64(plan.PriceUSD)

	// Node pricing
	cpuCost := float64(inst.CPU) * plan.NodePricing.CPUUSD
	memCost := float64(inst.MemoryMB/1024) * plan.NodePricing.MemoryUSD
	storageCost := float64(inst.StorageGB) * plan.NodePricing.StorageUSD

	total += cpuCost + memCost + storageCost

	// Engine multiplier
	engineMult := p.EnginePricing[inst.Engine]
	if engineMult == 0 {
		engineMult = 1.0
	}
	total *= engineMult

	// Region multiplier
	rp, ok := p.RegionPricing[inst.Region]
	if ok {
		total *= rp.Multiplier
	}

	// HA cost
	if plan.HA {
		total *= 1.25 // 25% HA surcharge
	}

	// Read replicas cost
	if plan.ReadReplicas > 0 {
		total += float64(plan.ReadReplicas) * (cpuCost + memCost) * 0.50
	}

	return total, nil
}
