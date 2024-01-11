# In memory cache

# Example of use:

	func main() {
		cache := NewCache(2)

		cache.Set("key1", "value1", 5*time.Second)

		go func() {
			time.Sleep(3 * time.Second)
			if err := cache.Update("key1", "updatedValue1"); err != nil {
				log.Println(err)
			}
		}()

		go func() {
			time.Sleep(7 * time.Second)
			if err := cache.Delete("key1"); err != nil {
				log.Println(err)
			}
		}()

		data, err := cache.Get("key1")
		if err != nil {
			log.Println(err)
		} else {
			log.Printf("Value for key1: %+v\n", data.items["key1"].value)
		}

		time.Sleep(10 * time.Second)
		cache.Close()
	}
