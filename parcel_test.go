package main

import (
	"database/sql"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// randSource is a source of random integers
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange uses randSource for generating random integers
	randRange = rand.New(randSource)
)

// getTestParcel returns test parcel
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete tests adding, getting and deleting
func TestAddGetDelete(t *testing.T) {
	// preparing connection
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err, "Require: Driver's error")
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// testing Add()
	number, err := store.Add(parcel)
	require.NoError(t, err, "Require: Add error")
	assert.NotEqual(t, 0, number, "Assert: wrong number returning")

	// testing Get()
	getParcel, err := store.Get(number)
	require.NoError(t, err, "Require: Get error")
	// comparing values in got parcel and test parcel
	assert.Equal(t, getParcel.Number, number, "Assert: wrong number")
	assert.Equal(t, getParcel.Client, parcel.Client, "Assert: wrong client")
	assert.Equal(t, getParcel.Address, parcel.Address, "Assert: wrong address")
	assert.Equal(t, getParcel.Status, parcel.Status, "Assert: wrong status")
	assert.Equal(t, getParcel.CreatedAt, parcel.CreatedAt, "Assert: wrong created_at")

	// testing Delete()
	err = store.Delete(number)
	require.NoError(t, err, "Require: Delete error")
	// checking if deleted parcel can't be got
	_, err = store.Get(number)
	require.Error(t, err, "Require: wrong parcel was deleted")
}

// TestSetAddress tests setting new address
func TestSetAddress(t *testing.T) {
	// preparing connection
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err, "Require: Driver's error")
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// adding new parcel
	number, err := store.Add(parcel)
	require.NoError(t, err, "Require: Add error")
	assert.NotEqual(t, 0, number, "Assert: wrong number returning")

	// seting new address
	newAddress := "new test address"
	err = store.SetAddress(number, newAddress)
	require.NoError(t, err, "Require: Set address error")

	// checking if address was changed
	getParcel, err := store.Get(number)
	require.NoError(t, err, "Require: Get error")
	// comparing current address and new address
	assert.Equal(t, getParcel.Address, newAddress, "Assert: adresses are not equal")
}

// TestSetStatus tests setting new status
func TestSetStatus(t *testing.T) {
	// preparing connection
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err, "Require: Driver's error")
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// adding new parcel
	number, err := store.Add(parcel)
	require.NoError(t, err, "Require: Add error")
	assert.NotNil(t, number, "Assert: number is NIL")

	// setting new status
	err = store.SetStatus(number, ParcelStatusSent)
	require.NoError(t, err, "Require: Set status error")

	// checking if status was changed
	getParcel, err := store.Get(number)
	require.NoError(t, err, "Require: Get error")
	// comparing current status and new status
	assert.Equal(t, getParcel.Status, ParcelStatusSent, "Assert: statuses are not equal")
}

// TestGetByClient tests parcels by client
func TestGetByClient(t *testing.T) {
	// preparing connection
	db, err := sql.Open("sqlite", "tracker.db")
	require.NoError(t, err, "Require: Driver's error")
	defer db.Close()

	store := NewParcelStore(db)
	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}
	parcelMap := map[int]Parcel{}

	// setting random value for each parcel
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	// adding parcels from slice into DB
	for i := 0; i < len(parcels); i++ {
		id, err := store.Add(parcels[i])
		require.NoError(t, err, "Require: Get error")
		assert.NotNil(t, id, "Assert: number is NIL")
		// setting new id for the parcel
		parcels[i].Number = id
		// adding parcel into map as a value with id as a key
		parcelMap[id] = parcels[i]
	}

	// getting list of parcels by client
	storedParcels, err := store.GetByClient(client)
	require.NoError(t, err, "Require: GetByClient error")
	// checking length of slice of parcels from DB and length of test slice of parcels
	assert.Len(t, storedParcels, len(parcels), "Assert: lemgths of parcels are not equal")

	// checking parcels from slice and map
	for _, parcel := range storedParcels {
		// checking if parcel from slice is in map
		assert.Contains(t, parcelMap, parcel.Number, "Assert: parcel is not in the map")
		// comparing parcels in slice and map
		assert.Equal(t, parcel, parcelMap[parcel.Number],
			"Assert: compared parcels are not equal")
	}
}
